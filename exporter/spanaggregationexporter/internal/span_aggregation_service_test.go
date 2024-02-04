package internal

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"go.uber.org/zap"
)

type DummySpanAggregationRepositoryImpl struct {
	logger   *zap.Logger
	toCreate []*ent.SpanAggregation
}

func CreateDummySpanAggregationRepositoryImpl() *DummySpanAggregationRepositoryImpl {
	return &DummySpanAggregationRepositoryImpl{logger: zap.NewExample()}
}

func (repository *DummySpanAggregationRepositoryImpl) FindAllByTraceIds(ctx context.Context, traceIds ...string) ([]*ent.SpanAggregation, error) {
	return []*ent.SpanAggregation{}, nil
}

func (repository *DummySpanAggregationRepositoryImpl) SaveAll(ctx context.Context, toCreate []*ent.SpanAggregation) error {
	repository.logger.Debug("save all", zap.Any("toCreate", toCreate))
	repository.toCreate = toCreate
	return nil
}

func (DummySpanAggregationRepositoryImpl) Shutdown(ctx context.Context) error {
	return nil
}

func TestSpanAggregationServiceImpl_Save(t *testing.T) {
	tests := []struct {
		name      string
		newSpans  []*ent.SpanAggregation
		results   map[string]struct{}
		validates func(*ent.SpanAggregation) bool
	}{
		{
			name:     "add the first non-root span of the trace",
			newSpans: []*ent.SpanAggregation{buildSpanAggregation(args{spanId: "test-first-non-root", parentSpanId: uuid.NewString()})},
			results: map[string]struct{}{
				"test-first-non-root": {},
			},
			validates: func(span *ent.SpanAggregation) bool { return true },
		},
		{
			name:     "add the first root span of the trace",
			newSpans: []*ent.SpanAggregation{buildSpanAggregation(args{spanId: "test-first-root"})},
			results: map[string]struct{}{
				"test-first-root": {},
			},
			validates: func(span *ent.SpanAggregation) bool { return true },
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
			results: map[string]struct{}{
				"test-parent": {},
				"test-child":  {},
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
			repository := CreateDummySpanAggregationRepositoryImpl()
			channel := make(chan rxgo.Item, 100)
			service := &SpanAggregationServiceImpl{
				repository:             repository,
				logger:                 zap.NewExample(),
				intervalInMilliseconds: 5,
				batchSize:              1000,
				batchChannel:           channel,
				observable:             rxgo.FromChannel(channel),
			}
			ctx, cancel := context.WithCancel(context.Background())
			service.Start(ctx)
			defer service.Shutdown(ctx)

			if err := service.Save(ctx, tt.newSpans); err != nil {
				t.Errorf("SpanAggregationServiceImpl.Save() error = %v", err)
			}
			time.Sleep(10 * time.Millisecond)
			for _, toCreate := range repository.toCreate {
				_, ok := tt.results[toCreate.ID]
				if !ok {
					t.Errorf("repository want %v", toCreate.ID)
				} else if !tt.validates(toCreate) {
					t.Errorf("repository get %s, but not pass the validation %s", toCreate.ID, toCreate)
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
