package internal

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
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

	intervalInMilliseconds int
	batchSize              int
	batchChannel           chan *channelEntry

	mutex *sync.Mutex
}

func CreateSpanAggregationServiceImpl(cacheConfig *CacheConfig, batchConfig *BatchConfig, repository SpanAggregationRepository, logger *zap.Logger) *SpanAggregationServiceImpl {
	traceCache := expirable.NewLRU[string, []*ent.SpanAggregation](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	parentCache := expirable.NewLRU[string, *ent.SpanAggregation](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	childrenCache := expirable.NewLRU[string, []*ent.SpanAggregation](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	mutext := &sync.Mutex{}
	return &SpanAggregationServiceImpl{
		logger, repository,
		traceCache, parentCache, childrenCache,
		batchConfig.BatchSize, batchConfig.IntervalInMilliseconds, make(chan *channelEntry, 100),
		mutext,
	}
}

func (service *SpanAggregationServiceImpl) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(service.intervalInMilliseconds))
	toCreateBatch := make([]*ent.SpanAggregation, 0, service.batchSize)
	toUpdateBatch := make([]*ent.SpanAggregation, 0, service.batchSize)
	// a new span is always to be created
	// so the lightweight delete operation can have a smaller where clause
	pending := make(map[string]int)

	flushBatchFunc := func(threshold int) {
		if len(toCreateBatch)+len(toUpdateBatch) > threshold {
			service.logger.Debug("flushing batch")
			err := service.repository.SaveAll(ctx, toCreateBatch, toUpdateBatch)
			if err != nil {
				service.logger.Error("failed to save batch", zap.Error(err))
			}
			for _, span := range toCreateBatch {
				delete(pending, span.ID)
			}
			toCreateBatch = toCreateBatch[:0]
			toUpdateBatch = toUpdateBatch[:0]
		}
	}

	go func(ctx context.Context) {
		service.logger.Debug("span aggregation service started")
		defer service.logger.Debug("span aggregation service stopped")
		for {
			select {
			case entry, ok := <-service.batchChannel:
				service.logger.Sugar().Debugf("receiving span %v %v %v", entry.span, entry.op, ok)
				if !ok {
					flushBatchFunc(0)
					break
				}

				if entry == nil {
					break
				}

				if entry.op == create {
					toCreateBatch = append(toCreateBatch, entry.span)
					pending[entry.span.ID] = 1
				} else if entry.op == update && pending[entry.span.ID] != 1 {
					toUpdateBatch = append(toUpdateBatch, entry.span)
				}
				flushBatchFunc(service.batchSize)
				ticker.Reset(time.Millisecond * time.Duration(service.intervalInMilliseconds))
			case <-ticker.C:
				service.logger.Debug("flushing batch by timer")
				flushBatchFunc(0)
			case <-ctx.Done():
				close(service.batchChannel)
				ticker.Stop()
				flushBatchFunc(0)
				return
			}
		}
	}(ctx)
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
		service.batchChannel <- entry
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
	// TODO use lock per trace
	service.mutex.Lock()
	defer service.mutex.Unlock()

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
	service.traceCache.Purge()
	service.parentCache.Purge()
	service.childrenCache.Purge()
	return service.repository.Shutdown(ctx)
}
