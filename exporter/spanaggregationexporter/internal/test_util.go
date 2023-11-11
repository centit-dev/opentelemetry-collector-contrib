package internal

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/schema"
)

type args struct {
	traceId      string
	spanId       string
	parentSpanId string
	duration     int64
	timestamp    time.Time
}

func (args *args) TraceId() string {
	if args.traceId == "" {
		return uuid.NewString()
	}
	return args.traceId
}

func (args *args) SpanId() string {
	if args.spanId == "" {
		return uuid.NewString()
	}
	return args.spanId
}

func (args *args) Timestamp() time.Time {
	if args.timestamp.IsZero() {
		return time.Now()
	}
	return args.timestamp
}

func (args *args) Duration() int64 {
	if args.duration == 0 {
		return int64(rand.Intn(100)) * 1_000_000
	}
	return args.duration
}

func buildSpanAggregation(args args) *ent.SpanAggregation {
	services := []string{"app", "db", "redis", "kafka", "es"}
	resourceAttributes := &schema.Attributes{}
	resourceAttributes.Add("resource", "value")
	spanAttributes := &schema.Attributes{}
	spanAttributes.Add("span", "value")
	return &ent.SpanAggregation{
		Timestamp:          args.Timestamp(),
		TraceId:            args.TraceId(),
		ParentSpanId:       args.parentSpanId,
		ID:                 args.SpanId(),
		PlatformName:       "beijing",
		RootServiceName:    "app",
		RootSpanName:       fmt.Sprintf("GET /api/%d", rand.Intn(10)),
		ServiceName:        services[rand.Intn(len(services))],
		SpanName:           fmt.Sprintf("GET /api/%d", rand.Intn(10)),
		ResourceAttributes: resourceAttributes,
		SpanAttributes:     spanAttributes,
		Duration:           args.Duration(),
		Gap:                0,
		SelfDuration:       0,
	}
}
