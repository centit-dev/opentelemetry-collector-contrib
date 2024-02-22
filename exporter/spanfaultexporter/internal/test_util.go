package internal

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
)

func buildSpanFault(traceId string) *spanTreeItem {
	fault := &ent.SpanFault{
		Timestamp:    time.Now(),
		ID:           traceId,
		PlatformName: []string{"online", "offline"}[rand.Intn(2)],
		AppCluster:   []string{"beijing", "shanghai"}[rand.Intn(2)],
		InstanceName: []string{"instance1", "instance2"}[rand.Intn(2)],
		ParentSpanId: uuid.NewString(),
		SpanId:       uuid.NewString(),
		ServiceName:  []string{"app1", "app2"}[rand.Intn(2)],
		SpanName:     []string{"span1", "span2"}[rand.Intn(2)],
		FaultKind:    []string{"Business", "System"}[rand.Intn(2)],
	}
	return &spanTreeItem{fault, rand.Int63n(1000), nil, 0}
}
