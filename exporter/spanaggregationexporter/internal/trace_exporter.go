package internal

import (
	"context"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/schema"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"
	"go.uber.org/zap"
)

const platFormNameKey = "k8s.pod.platform"

type TraceExporter struct {
	logger  *zap.Logger
	service SpanAggregationService
}

func CreateTraceExporter(service SpanAggregationService, logger *zap.Logger) *TraceExporter {
	return &TraceExporter{
		logger:  logger,
		service: service,
	}
}

func (t *TraceExporter) Start(ctx context.Context, _ component.Host) error {
	t.service.Start(ctx)
	return nil
}

func (t *TraceExporter) PushTraceData(ctx context.Context, traces ptrace.Traces) error {
	slice := traces.ResourceSpans()
	for i := 0; i < slice.Len(); i++ {
		resource := slice.At(i)
		resourceAttributes := resource.Resource().Attributes()
		resourceAttributeMap := attributesToMap(resourceAttributes)
		var platformName string
		if platformNameValue, exists := resourceAttributes.Get(platFormNameKey); exists {
			platformName = platformNameValue.Str()
		}
		var serviceName string
		if serviceNameValue, exists := resourceAttributes.Get(conventions.AttributeServiceName); exists {
			serviceName = serviceNameValue.Str()
		}

		batches := resource.ScopeSpans()
		for j := 0; j < batches.Len(); j++ {
			scopeSpans := batches.At(j)
			batch := scopeSpans.Spans()
			for k := 0; k < batch.Len(); k++ {
				span := batch.At(k)
				aggregation := t.buildSpanAggregation(resourceAttributeMap, &span, platformName, serviceName)
				err := t.service.Save(ctx, aggregation)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *TraceExporter) buildSpanAggregation(resourceAttributeMap *schema.Attributes, span *ptrace.Span, platformName string, serviceName string) *ent.SpanAggregation {
	duration := span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).Nanoseconds()
	spanAggregation := &ent.SpanAggregation{
		Timestamp:          span.StartTimestamp().AsTime(),
		TraceId:            span.TraceID().String(),
		ParentSpanId:       span.ParentSpanID().String(),
		ID:                 span.SpanID().String(),
		PlatformName:       platformName,
		ServiceName:        serviceName,
		SpanName:           span.Name(),
		ResourceAttributes: resourceAttributeMap,
		SpanAttributes:     attributesToMap(span.Attributes()),
		Duration:           duration,
		SelfDuration:       duration,
	}
	return spanAggregation
}

func attributesToMap(attributes pcommon.Map) *schema.Attributes {
	m := &schema.Attributes{}
	attributes.Range(func(k string, v pcommon.Value) bool {
		m.Add(k, v.AsString())
		return true
	})
	return m
}

func (t *TraceExporter) Shutdown(ctx context.Context) error {
	return t.service.Shutdown(ctx)
}
