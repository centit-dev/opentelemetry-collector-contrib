package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"go.uber.org/zap"
)

// disabled for integration test
func DisabledTestSpanAggregationRepositoryImpl_SaveAll(t *testing.T) {
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
		updated, err := repository.FindAllByTraceIds(context.Background(), traceId)
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

// disabled for integration test
func DisabledTestSpanAggregationRepositoryImpl_formatMap(t *testing.T) {
	attributes := map[string]string{
		"db.connection_string": "mysql://smartobserv-mysql.smartobserv-local:3306",
		"db.user":              "root",
		"thread.name":          "http-nio-8080-exec-8",
		"db.statement.values":  "['max_minutes_remain_topay']",
		"db.name":              "starshop_test",
		"db.statement":         "select bizparamet0_.code as code1_6_0_, bizparamet0_.value as value2_6_0_ from tb_sys_parameter bizparamet0_ where bizparamet0_.code=?",
		"db.system":            "mysql",
		"net.peer.port":        "3306",
		"net.peer.name":        "smartobserv-mysql.smartobserv-local",
		"db.sql.table":         "tb_sys_parameter",
		"db.operation":         "SELECT",
		"thread.id":            "44",
	}

	for key, value := range attributes {
		span := buildSpanAggregation(args{})
		span.SpanAttributes.Add(key, value)
		logger := zap.NewExample()
		client, _ := CreateClient(&ClickHouseConfig{
			Endpoint: "host.docker.internal:29000",
			Database: "otel",
			Debug:    true,
		}, logger)
		repository := &SpanAggregationRepositoryImpl{client: client, logger: logger}

		t.Run(fmt.Sprintf("format %s", key), func(t *testing.T) {
			result := formatMap(span.SpanAttributes)
			logger.Info(result)
			err := repository.SaveAll(context.Background(), []*ent.SpanAggregation{span}, nil)
			if err != nil {
				t.Errorf("SpanAggregationRepositoryImpl.SaveAll() error = %v", err)
				return
			}
		})
	}
}
