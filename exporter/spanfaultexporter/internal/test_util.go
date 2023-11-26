package internal

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
)

func buildSpanFault(traceId string) *ent.SpanFault {
	return &ent.SpanFault{
		Timestamp:    time.Now(),
		TraceId:      traceId,
		PlatformName: []string{"online", "offline"}[rand.Intn(2)],
		ClusterName:  []string{"beijing", "shanghai"}[rand.Intn(2)],
		InstanceName: []string{"instance1", "instance2"}[rand.Intn(2)],
		ParentSpanId: uuid.NewString(),
		ID:           uuid.NewString(),
		ServiceName:  []string{"app1", "app2"}[rand.Intn(2)],
		SpanName:     []string{"span1", "span2"}[rand.Intn(2)],
		FaultKind:    []string{"Business", "System"}[rand.Intn(2)],
		IsRoot:       []bool{true, false}[rand.Intn(2)],
	}
}
