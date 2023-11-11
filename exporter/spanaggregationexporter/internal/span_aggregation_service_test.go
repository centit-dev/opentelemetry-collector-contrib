package internal

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"go.uber.org/zap"
)

type DummySpanAggregationRepositoryImpl struct {
	logger   *zap.Logger
	spans    []*ent.SpanAggregation
	toCreate []*ent.SpanAggregation
	toUpdate []*ent.SpanAggregation
}

func CreateDummySpanAggregationRepositoryImpl(spans []*ent.SpanAggregation) *DummySpanAggregationRepositoryImpl {
	return &DummySpanAggregationRepositoryImpl{logger: zap.NewExample(), spans: spans}
}

func (repository *DummySpanAggregationRepositoryImpl) FindAllByTraceId(ctx context.Context, traceId string) ([]*ent.SpanAggregation, error) {
	return repository.spans, nil
}

func (repository *DummySpanAggregationRepositoryImpl) SaveAll(ctx context.Context, toCreate []*ent.SpanAggregation, toUpdate []*ent.SpanAggregation) error {
	repository.logger.Debug("save all", zap.Any("toCreate", toCreate), zap.Any("toUpdate", toUpdate))
	repository.toCreate = toCreate
	repository.toUpdate = toUpdate
	return nil
}

func (DummySpanAggregationRepositoryImpl) Shutdown(ctx context.Context) error {
	return nil
}

func TestSpanAggregationServiceImpl_Save(t *testing.T) {
	tests := []struct {
		name      string
		newSpans  []*ent.SpanAggregation
		cached    []*ent.SpanAggregation
		results   map[string]channelEntryOp
		validates func(*ent.SpanAggregation) bool
	}{
		{
			name:     "add the first non-root span of the trace",
			newSpans: []*ent.SpanAggregation{buildSpanAggregation(args{spanId: "test-first-non-root", parentSpanId: uuid.NewString()})},
			cached:   []*ent.SpanAggregation{},
			results: map[string]channelEntryOp{
				"test-first-non-root": create,
			},
			validates: func(span *ent.SpanAggregation) bool { return true },
		},
		{
			name:     "add the first root span of the trace",
			newSpans: []*ent.SpanAggregation{buildSpanAggregation(args{spanId: "test-first-root"})},
			cached:   []*ent.SpanAggregation{},
			results: map[string]channelEntryOp{
				"test-first-root": create,
			},
			validates: func(span *ent.SpanAggregation) bool { return true },
		},

		{
			name: "add a new non-root span to the trace which has no root",
			newSpans: []*ent.SpanAggregation{
				buildSpanAggregation(args{spanId: "test-follow-non-root", traceId: "shared-trace-id", parentSpanId: uuid.NewString()}),
			},
			cached: []*ent.SpanAggregation{
				buildSpanAggregation(args{traceId: "shared-trace-id", parentSpanId: uuid.NewString()}),
			},
			results: map[string]channelEntryOp{
				"test-follow-non-root": create,
			},
			validates: func(span *ent.SpanAggregation) bool { return true },
		},
		{
			name: "add a new non-root span to the trace which has a root",
			newSpans: []*ent.SpanAggregation{
				buildSpanAggregation(args{spanId: "test-follow-non-root", traceId: "shared-trace-id", parentSpanId: uuid.NewString()}),
			},
			cached: []*ent.SpanAggregation{
				buildSpanAggregation(args{traceId: "shared-trace-id"}),
			},
			results: map[string]channelEntryOp{
				"test-follow-non-root": create,
			},
			validates: func(span *ent.SpanAggregation) bool { return span.RootServiceName != "" },
		},
		{
			name: "add a new root span to the trace which has no root",
			newSpans: []*ent.SpanAggregation{
				buildSpanAggregation(args{spanId: "test-follow-non-root", traceId: "shared-trace-id"}),
			},
			cached: []*ent.SpanAggregation{
				buildSpanAggregation(args{spanId: "update-root-value", traceId: "shared-trace-id", parentSpanId: uuid.NewString()}),
			},
			results: map[string]channelEntryOp{
				"test-follow-non-root": create,
				"update-root-value":    update,
			},
			validates: func(span *ent.SpanAggregation) bool { return span.RootServiceName != "" },
		},

		{
			name: "add a new span father to one of the spans",
			newSpans: []*ent.SpanAggregation{
				buildSpanAggregation(args{
					spanId:    "test-parent",
					traceId:   "shared-trace-id",
					timestamp: time.Unix(0, 0),
				}),
			},
			cached: []*ent.SpanAggregation{
				buildSpanAggregation(args{
					spanId:       "test-child",
					traceId:      "shared-trace-id",
					parentSpanId: "test-parent",
					timestamp:    time.Unix(0, 0).Add(10 * time.Nanosecond),
				}),
			},
			results: map[string]channelEntryOp{
				"test-parent": create,
				"test-child":  update,
			},
			validates: func(span *ent.SpanAggregation) bool { return span.ID != "test-child" || span.Gap == 10 },
		},
		{
			name: "add a new span child to one of the spans",
			newSpans: []*ent.SpanAggregation{
				buildSpanAggregation(args{
					spanId:       "test-child",
					traceId:      "shared-trace-id",
					parentSpanId: "test-parent",
					timestamp:    time.Unix(0, 0).Add(10 * time.Nanosecond),
				}),
			},
			cached: []*ent.SpanAggregation{
				buildSpanAggregation(args{
					spanId:    "test-parent",
					traceId:   "shared-trace-id",
					timestamp: time.Unix(0, 0),
				}),
			},
			results: map[string]channelEntryOp{
				"test-child":  create,
				"test-parent": update,
			},
			validates: func(span *ent.SpanAggregation) bool { return span.ID != "test-child" || span.Gap == 10 },
		},

		{
			name: "add 2 new spans and 1 of them is the father of the other",
			newSpans: []*ent.SpanAggregation{
				buildSpanAggregation(args{
					spanId:    "test-parent",
					traceId:   "shared-trace-id",
					timestamp: time.Unix(0, 0),
					duration:  30,
				}),
				buildSpanAggregation(args{
					spanId:       "test-child",
					traceId:      "shared-trace-id",
					parentSpanId: "test-parent",
					timestamp:    time.Unix(0, 0).Add(10 * time.Nanosecond),
					duration:     10,
				}),
			},
			cached: []*ent.SpanAggregation{},
			results: map[string]channelEntryOp{
				"test-parent": create,
				"test-child":  create,
			},
			validates: func(span *ent.SpanAggregation) bool {
				if span.ID == "test-parent" {
					return span.Gap == 0 && span.SelfDuration == 20
				}
				if span.ID == "test-child" {
					return span.Gap == 10
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := CreateDummySpanAggregationRepositoryImpl(tt.cached)
			service := &SpanAggregationServiceImpl{
				repository:             repository,
				logger:                 zap.NewExample(),
				traceCache:             expirable.NewLRU[string, []*ent.SpanAggregation](1000, nil, time.Minute),
				parentCache:            expirable.NewLRU[string, *ent.SpanAggregation](1000, nil, time.Minute),
				childrenCache:          expirable.NewLRU[string, []*ent.SpanAggregation](1000, nil, time.Minute),
				intervalInMilliseconds: 5,
				batchSize:              1000,
				batchChannel:           make(chan *channelEntry, 100),
				mutex:                  &sync.Mutex{},
			}
			ctx, cancel := context.WithCancel(context.Background())
			service.Start(ctx)
			defer service.Shutdown(ctx)

			for _, span := range tt.newSpans {
				if err := service.Save(ctx, span); err != nil {
					t.Errorf("SpanAggregationServiceImpl.Save() error = %v for %v", err, span)
				}
			}
			time.Sleep(10 * time.Millisecond)
			if len(repository.toCreate)+len(repository.toUpdate) != len(tt.results) {
				t.Errorf("repository want %d, get %d", len(tt.results), len(repository.toCreate)+len(repository.toUpdate))
			}
			for _, toCreate := range repository.toCreate {
				op, ok := tt.results[toCreate.ID]
				if ok && op != create {
					t.Errorf("repository want %s %v, get %v", toCreate.ID, op, create)
				} else if !tt.validates(toCreate) {
					t.Errorf("repository get %s, but not pass the validation %s", toCreate.ID, toCreate)
				}
			}
			for _, toUpdate := range repository.toUpdate {
				op, ok := tt.results[toUpdate.ID]
				if ok && op != update {
					t.Errorf("repository want %s %v, get %v", toUpdate.ID, op, update)
				} else if !tt.validates(toUpdate) {
					t.Errorf("repository get %s, but not pass the validation %s", toUpdate.ID, toUpdate)
				}
			}
			cancel()
		})
	}
}

func TestSpanAggregationServiceImpl_calculateSelfDuration(t *testing.T) {
	tests := []struct {
		name            string
		spanDuration    int64
		childrenOffsets []int64
		want            int64
	}{
		{
			name:            "Children not overlap",
			spanDuration:    100,
			childrenOffsets: []int64{0, 10, 20, 30, 40, 50},
			want:            70,
		},
		{
			name:            "Children overlap",
			spanDuration:    100,
			childrenOffsets: []int64{0, 15, 10, 30, 40, 50},
			want:            60,
		},
		{
			name:            "Children overlap with wrong order",
			spanDuration:    100,
			childrenOffsets: []int64{10, 30, 0, 15, 40, 50},
			want:            60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &SpanAggregationServiceImpl{}
			span, children := buildDurationTestCase(tt.spanDuration, tt.childrenOffsets...)
			if got := service.calculateSelfDuration(span, children); got != tt.want {
				t.Errorf("SpanAggregationServiceImpl.calculateSelfDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func buildDurationTestCase(duration int64, timeOffsets ...int64) (*ent.SpanAggregation, []*ent.SpanAggregation) {
	now := time.Unix(0, 0)
	span := buildSpanAggregation(args{timestamp: now, duration: duration})
	children := make([]*ent.SpanAggregation, len(timeOffsets)/2)
	for index := 0; index < len(timeOffsets); index += 2 {
		start := now.Add(time.Duration(timeOffsets[index]) * time.Nanosecond)
		end := now.Add(time.Duration(timeOffsets[index+1]) * time.Nanosecond)
		children[index/2] = buildSpanAggregation(args{timestamp: start, duration: end.Sub(start).Nanoseconds()})
	}
	return span, children
}
