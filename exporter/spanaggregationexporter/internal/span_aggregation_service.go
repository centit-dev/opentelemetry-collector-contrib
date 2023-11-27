package internal

import (
	"context"
	"sort"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"go.uber.org/zap"
)

type CacheConfig struct {
	MaxSize         int `mapstructure:"max_size"`
	ExpireInMinutes int `mapstructure:"expire_in_minutes"`
}

type BatchConfig struct {
	BatchSize int `mapstructure:"batch_size"`
	// use millisecond for testing purpose and be more flexible
	IntervalInMilliseconds int `mapstructure:"interval_in_milliseconds"`
}

type channelEntryOp int

const (
	create channelEntryOp = iota
	update
)

type channelEntry struct {
	op   channelEntryOp
	span *ent.SpanAggregation
}

type SpanAggregationService interface {
	Start(ctx context.Context)
	// extract all spans from a batch
	// and get cache or persistent copies from the database if the cache misses
	// group spans by span id, parent id, trace id
	// and get root resource attributes from the trace id group
	// and get gap from the parent group
	// and get self duration from the span id group
	// then build the span aggregation for creating or updating
	// and update the cache and the database
	Save(ctx context.Context, span *ent.SpanAggregation) error
	Shutdown(ctx context.Context) error
}

type SpanAggregationServiceImpl struct {
	logger     *zap.Logger
	repository SpanAggregationRepository

	traceCache    *expirable.LRU[string, []*ent.SpanAggregation]
	parentCache   *expirable.LRU[string, *ent.SpanAggregation]
	childrenCache *expirable.LRU[string, []*ent.SpanAggregation]

	batchChannel           chan rxgo.Item
	observable             rxgo.Observable
	intervalInMilliseconds int
	batchSize              int
}

func CreateSpanAggregationServiceImpl(cacheConfig *CacheConfig, batchConfig *BatchConfig, repository SpanAggregationRepository, logger *zap.Logger) *SpanAggregationServiceImpl {
	traceCache := expirable.NewLRU[string, []*ent.SpanAggregation](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	parentCache := expirable.NewLRU[string, *ent.SpanAggregation](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	childrenCache := expirable.NewLRU[string, []*ent.SpanAggregation](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	channel := make(chan rxgo.Item, 1000)
	return &SpanAggregationServiceImpl{
		logger:                 logger,
		repository:             repository,
		traceCache:             traceCache,
		parentCache:            parentCache,
		childrenCache:          childrenCache,
		batchChannel:           channel,
		observable:             rxgo.FromChannel(channel),
		intervalInMilliseconds: batchConfig.IntervalInMilliseconds,
		batchSize:              batchConfig.BatchSize,
	}
}

func (service *SpanAggregationServiceImpl) Start(ctx context.Context) {
	service.observable.
		BufferWithTimeOrCount(
			rxgo.WithDuration(time.Millisecond*time.Duration(service.intervalInMilliseconds)),
			service.batchSize).
		DoOnNext(func(items interface{}) {
			service.logger.Debug("received batch")
			values, ok := items.([]interface{})
			if !ok {
				return
			}
			seen := make(map[string]bool, len(values))
			creates := make([]*ent.SpanAggregation, 0, len(values))
			updates := make([]*ent.SpanAggregation, 0, len(values))
			for _, item := range values {
				entry, ok := item.(*channelEntry)
				if !ok {
					continue
				}
				if seen[entry.span.ID] {
					continue
				}
				if entry.op == create {
					creates = append(creates, entry.span)
				} else if entry.op == update {
					updates = append(updates, entry.span)
				}
				seen[entry.span.ID] = true
				err := service.repository.SaveAll(ctx, creates, updates)
				if err != nil {
					service.logger.Error("failed to save batch", zap.Error(err))
				}
			}
		}, rxgo.WithPool(10))
}

// the span is always new to the trace
func (service *SpanAggregationServiceImpl) Save(ctx context.Context, span *ent.SpanAggregation) error {
	upsertBatch := make(map[string]*channelEntry)

	err := service.buildCacheByTraceId(ctx, span.TraceId)
	if err != nil {
		return err
	}
	service.updateTraceGroup(ctx, &upsertBatch, span)
	service.updateParentGroup(ctx, &upsertBatch, span)
	service.updateChildGroup(ctx, &upsertBatch, span)

	for _, entry := range upsertBatch {
		service.logger.Sugar().Debugf("saving span %v %v", entry.span, entry.op)
		service.batchChannel <- rxgo.Of(entry)
	}
	return nil
}

// trace group:
func (service *SpanAggregationServiceImpl) updateTraceGroup(ctx context.Context, updateBatch *map[string]*channelEntry, newSpan *ent.SpanAggregation) {
	spans, exists := service.traceCache.Get(newSpan.TraceId)
	// add span to trace group
	if exists {
		spans = append(spans, newSpan)
	} else {
		spans = []*ent.SpanAggregation{newSpan}
	}
	service.traceCache.Add(newSpan.TraceId, spans)

	var root *ent.SpanAggregation
	for _, span := range spans {
		if span.ParentSpanId == "" {
			root = span
		}
	}

	if root != nil {
		if root.ID == newSpan.ID {
			// if it's the root, update root attributes for all
			for _, span := range spans {
				span.RootServiceName = root.ServiceName
				span.RootSpanName = root.SpanName
				op := update
				if span.ID == newSpan.ID {
					op = create
				}
				(*updateBatch)[span.ID] = &channelEntry{op, span}
			}
		} else {
			// else if the root is in there, update root attributes for it
			newSpan.RootServiceName = root.ServiceName
			newSpan.RootSpanName = root.SpanName
		}
	}
	(*updateBatch)[newSpan.ID] = &channelEntry{op: create, span: newSpan}
}

// parent group:
func (service *SpanAggregationServiceImpl) updateParentGroup(ctx context.Context, updateBatch *map[string]*channelEntry, newSpan *ent.SpanAggregation) {
	// add span to siblings group as a child
	siblings, exists := service.childrenCache.Get(newSpan.ParentSpanId)
	if exists {
		siblings = append(siblings, newSpan)
	} else {
		siblings = []*ent.SpanAggregation{newSpan}
	}
	service.childrenCache.Add(newSpan.ParentSpanId, siblings)

	// update its parent's self duration
	if parent, exists := service.parentCache.Get(newSpan.ParentSpanId); exists {
		parent.SelfDuration = service.calculateSelfDuration(parent, siblings)
		(*updateBatch)[parent.ID] = &channelEntry{op: update, span: parent}

		newSpan.Gap = newSpan.Timestamp.Sub(parent.Timestamp).Nanoseconds()
		(*updateBatch)[newSpan.ID] = &channelEntry{op: create, span: newSpan}
	}
}

// children group:
func (service *SpanAggregationServiceImpl) updateChildGroup(ctx context.Context, updateBatch *map[string]*channelEntry, newSpan *ent.SpanAggregation) {
	// add span to parent group as a parent
	service.parentCache.Add(newSpan.ID, newSpan)

	// if the span has children, update gap for them and update self duration for the span
	if children, exists := service.childrenCache.Get(newSpan.ID); exists {
		for _, child := range children {
			child.Gap = child.Timestamp.Sub(newSpan.Timestamp).Nanoseconds()
			(*updateBatch)[child.ID] = &channelEntry{op: update, span: child}
		}
		newSpan.SelfDuration = service.calculateSelfDuration(newSpan, children)
		(*updateBatch)[newSpan.ID] = &channelEntry{op: create, span: newSpan}
	}
}

// build cache by trace id
// update the local cache when new span is successfully saved
// so we don't have to query the database again
func (service *SpanAggregationServiceImpl) buildCacheByTraceId(ctx context.Context, traceId string) error {
	if service.traceCache.Contains(traceId) {
		return nil
	}

	aggregations, err := service.repository.FindAllByTraceId(ctx, traceId)
	if err != nil {
		return err
	}

	for _, aggregation := range aggregations {
		// cache span by trace id
		if spans, exists := service.traceCache.Get(aggregation.TraceId); exists {
			service.traceCache.Add(aggregation.TraceId, append(spans, aggregation))
		} else {
			service.traceCache.Add(aggregation.TraceId, []*ent.SpanAggregation{aggregation})
		}

		service.parentCache.Add(aggregation.ID, aggregation)

		if aggregation.ParentSpanId != "" {
			if children, exists := service.childrenCache.Get(aggregation.ParentSpanId); exists {
				service.childrenCache.Add(aggregation.ParentSpanId, append(children, aggregation))
			} else {
				service.childrenCache.Add(aggregation.ParentSpanId, []*ent.SpanAggregation{aggregation})
			}
		}
	}
	return nil
}

func (SpanAggregationServiceImpl) calculateSelfDuration(span *ent.SpanAggregation, children []*ent.SpanAggregation) int64 {
	if len(children) == 0 {
		return span.Duration
	}

	sort.Slice(children, func(i, j int) bool {
		return children[i].Timestamp.Before(children[j].Timestamp)
	})
	intervals := make([][]int64, 0, len(children))
	for _, child := range children {
		if len(intervals) == 0 {
			intervals = append(intervals, []int64{child.Timestamp.UnixNano(), child.Timestamp.UnixNano() + child.Duration, child.Duration})
			continue
		}

		last_entry := intervals[len(intervals)-1]
		last_start := last_entry[0]
		last_end := last_entry[1]

		current_end := child.Timestamp.UnixNano() + child.Duration
		if child.Timestamp.UnixNano() > last_end {
			// not overlap:
			// last:  |---------|
			// next:                |---------|
			intervals = append(intervals, []int64{child.Timestamp.UnixNano(), current_end, child.Duration})
		} else if current_end > last_end {
			// overlaps:
			// last:  |---------|
			// next:       |---------|
			intervals[len(intervals)-1] = []int64{last_start, current_end, current_end - last_start}
		}
		// sorting by start time acceding ganrantees start_time > last_start_time
		// so this is impossible
		// last:         |---------|
		// next: |---------|
	}
	selfDuration := span.Duration
	for _, interval := range intervals {
		selfDuration -= interval[2]
	}
	// client error may cause negative self duration
	if selfDuration < 0 {
		selfDuration = 0
	}
	return selfDuration
}

func (service *SpanAggregationServiceImpl) Shutdown(ctx context.Context) error {
	_, cancel := service.observable.Connect(ctx)
	cancel()
	close(service.batchChannel)
	service.traceCache.Purge()
	service.parentCache.Purge()
	service.childrenCache.Purge()
	return service.repository.Shutdown(ctx)
}
