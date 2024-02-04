package internal

import (
	"context"
	"sort"
	"time"

	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"go.uber.org/zap"
)

type BatchConfig struct {
	BatchSize int `mapstructure:"batch_size"`
	// use millisecond for testing purpose and be more flexible
	IntervalInMilliseconds int `mapstructure:"interval_in_milliseconds"`
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
	Save(ctx context.Context, spans []*ent.SpanAggregation) error
	Shutdown(ctx context.Context) error
}

type SpanAggregationServiceImpl struct {
	logger     *zap.Logger
	repository SpanAggregationRepository

	batchChannel           chan rxgo.Item
	observable             rxgo.Observable
	intervalInMilliseconds int
	batchSize              int
}

func CreateSpanAggregationServiceImpl(batchConfig *BatchConfig, repository SpanAggregationRepository, logger *zap.Logger) *SpanAggregationServiceImpl {
	channel := make(chan rxgo.Item, 1000)
	return &SpanAggregationServiceImpl{
		logger:                 logger,
		repository:             repository,
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
			creates := make([]*ent.SpanAggregation, 0, len(values))
			for _, item := range values {
				span, ok := item.(*ent.SpanAggregation)
				if !ok {
					continue
				}
				creates = append(creates, span)
				service.logger.Sugar().Debugf("saving span %v %v %v", span.TraceId, span.ID)
			}
			err := service.repository.SaveAll(ctx, creates)
			if err != nil {
				service.logger.Error("failed to save batch", zap.Error(err))
			}
		}, rxgo.WithPool(10))
}

// the span is always new to the trace
func (service *SpanAggregationServiceImpl) Save(ctx context.Context, aggregations []*ent.SpanAggregation) error {
	traceRoots := make(map[string]*ent.SpanAggregation)
	traces := make(map[string][]*ent.SpanAggregation)
	parents := make(map[string]*ent.SpanAggregation)
	childrenGroup := make(map[string][]*ent.SpanAggregation)

	for _, aggregation := range aggregations {
		if aggregation.ParentSpanId == "" {
			traceRoots[aggregation.TraceId] = aggregation
		}
		// group span by trace id
		if spans, exists := traces[aggregation.TraceId]; exists {
			traces[aggregation.TraceId] = append(spans, aggregation)
		} else {
			traces[aggregation.TraceId] = []*ent.SpanAggregation{aggregation}
		}

		parents[aggregation.ID] = aggregation

		if aggregation.ParentSpanId != "" {
			if siblings, exists := childrenGroup[aggregation.ParentSpanId]; exists {
				childrenGroup[aggregation.ParentSpanId] = append(siblings, aggregation)
			} else {
				childrenGroup[aggregation.ParentSpanId] = []*ent.SpanAggregation{aggregation}
			}
		}
	}

	upsertBatch := make(map[string]*ent.SpanAggregation)
	// update roots
	for traceId, root := range traceRoots {
		spans := traces[traceId]
		for _, span := range spans {
			span.RootServiceName = root.ServiceName
			span.RootSpanName = root.SpanName
			upsertBatch[span.ID] = span
		}
	}
	// update gap and self duration
	for parentId, parent := range parents {
		upsertBatch[parent.ID] = parent
		if children, exists := childrenGroup[parentId]; exists {
			for _, child := range children {
				child.Gap = child.Timestamp.Sub(parent.Timestamp).Nanoseconds()
			}
			parent.SelfDuration = service.calculateSelfDuration(parent, children)
		} else {
			parent.SelfDuration = parent.Duration
		}
	}

	for _, entry := range upsertBatch {
		service.batchChannel <- rxgo.Of(entry)
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
	maxEndTime := span.Timestamp.Add(time.Duration(span.Duration))
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
		// sorting by start time ascending ganrantees start_time > last_start_time
		// so this is impossible
		// last:         |---------|
		// next: |---------|

		// lastly, if the last end time is greater than the parent end time, we need to adjust the last interval
		// last:            |---------|
		// parent: |---------|
		last_entry = intervals[len(intervals)-1]
		last_start = last_entry[0]
		last_end = last_entry[1]
		if last_end > maxEndTime.UnixNano() {
			intervals[len(intervals)-1] = []int64{last_start, maxEndTime.UnixNano(), maxEndTime.UnixNano() - last_start}
			// and there is no point to continue
			break
		}
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
	return service.repository.Shutdown(ctx)
}
