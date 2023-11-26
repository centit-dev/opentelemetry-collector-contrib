package internal

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
	"go.uber.org/zap"
)

func DisabledTestSpanFaultRepositoryImpl_SaveAll(t *testing.T) {
	logger := zap.NewExample()
	client, _ := CreateClient(&ClickHouseConfig{
		Endpoint: "host.docker.internal:29000",
		Database: "otel",
		Debug:    true,
	}, logger)
	repo := CreateSpanFaultRepository(client, logger)

	causes := make([]*ent.SpanFault, 0, 10)
	traceId := uuid.NewString()
	for i := 0; i < 10; i++ {
		causes = append(causes, buildSpanFault(traceId))
	}

	err := repo.SaveAll(context.Background(), causes, nil)
	if err != nil {
		t.Errorf("failed to save span faults: %v", err)
	}

	// update the first 5 of them
	for i := 0; i < 5; i++ {
		causes[i].FaultKind = ""
	}
	err = repo.SaveAll(context.Background(), nil, causes[:5])
	if err != nil {
		t.Errorf("failed to update span faults: %v", err)
	}

	// and they are updated
	causes, err = repo.GetSpanFaultsByTraceId(context.Background(), traceId)
	if err != nil {
		t.Errorf("failed to get span faults: %v", err)
	}
	for i := 0; i < 5; i++ {
		if causes[i].FaultKind != "" {
			t.Errorf("failed to update span faults: %v", err)
		}
	}
}
