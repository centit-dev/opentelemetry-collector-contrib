package internal

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"go.uber.org/zap"
)

// disabled for integration test
func TestSpanAggregationRepositoryImpl_SaveAll(t *testing.T) {
	logger := zap.NewExample()
	client, _ := CreateClient(&ClickHouseConfig{
		Endpoint: "host.docker.internal:29000",
		Database: "otel",
		Debug:    true,
	}, logger)
	repository := &SpanAggregationRepositoryImpl{client: client, logger: logger}

	t.Run("TestSaveAll", func(t *testing.T) {
		// create 10 aggregations
		aggregations := make([]*ent.SpanAggregation, 10)
		traceId := uuid.NewString()
		for i := 0; i < 10; i++ {
			aggregations[i] = buildSpanAggregation(args{traceId: traceId})
		}
		err := repository.SaveAll(context.Background(), aggregations, nil)
		if err != nil {
			t.Errorf("SpanAggregationRepositoryImpl.SaveAll() error = %v", err)
			return
		}

		// update 10 aggregations
		for _, aggregation := range aggregations {
			aggregation.SelfDuration = aggregation.SelfDuration + 1
		}
		err = repository.SaveAll(context.Background(), nil, aggregations)
		if err != nil {
			t.Errorf("SpanAggregationRepositoryImpl.SaveAll() error = %v", err)
			return
		}

		// verify 10 aggregations are saved
		updated, err := repository.FindAllByTraceId(context.Background(), traceId)
		if err != nil {
			t.Errorf("SpanAggregationRepositoryImpl.FindAllByTraceId() error = %v", err)
			return
		}
		if len(updated) != 10 {
			t.Errorf("SpanAggregationRepositoryImpl.FindAllByTraceId() = %v", updated)
			return
		}
		for _, aggregation := range updated {
			if aggregation.SelfDuration != 1 {
				t.Errorf("SpanAggregationRepositoryImpl.FindAllByTraceId() = %v", aggregation)
				return
			}
		}
	})
}
