package internal

import (
	"context"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent/schema"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"
	"go.uber.org/zap"
)

const (
	platFormNameKey = "service.platform"
	faultKindKey    = "fault.kind"
)

type TraceExporter struct {
	logger  *zap.Logger
	service SpanFaultService
}

func CreateTraceExporter(service SpanFaultService, logger *zap.Logger) *TraceExporter {
	return &TraceExporter{logger, service}
}

func (t *TraceExporter) Start(ctx context.Context, _ component.Host) error {
	t.service.Start(ctx)
	return nil
}

func (t *TraceExporter) PushTraceData(ctx context.Context, traces ptrace.Traces) error {
	items := make([]*spanTreeItem, 0, traces.SpanCount())
	slice := traces.ResourceSpans()
	for i := 0; i < slice.Len(); i++ {
		resource := slice.At(i)
		resourceAttributes := resource.Resource().Attributes()
		var platformName string
		if platformNameValue, exists := resourceAttributes.Get(platFormNameKey); exists {
			platformName = platformNameValue.Str()
		}
		var appCluster string
		if appClusterValue, exists := resourceAttributes.Get(conventions.AttributeK8SDeploymentName); exists {
			appCluster = appClusterValue.Str()
		}
		var instanceName string
		if instanceNameValue, exists := resourceAttributes.Get(conventions.AttributeK8SPodName); exists {
			instanceName = instanceNameValue.Str()
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
				item := t.buildFaultTreeItem(&resourceAttributes, &span, platformName, appCluster, instanceName, serviceName)
				items = append(items, item)
			}
		}
	}

	return t.service.Save(ctx, items)
}

func (t *TraceExporter) buildFaultTreeItem(attributes *pcommon.Map, span *ptrace.Span, platformName string, appCluster string, instanceName string, serviceName string) *spanTreeItem {
	var faultKind string
	if faultKindValue, exists := span.Attributes().Get(faultKindKey); exists {
		faultKind = faultKindValue.Str()
	}
	fault := &ent.SpanFault{
		Timestamp:          span.StartTimestamp().AsTime(),
		ID:                 span.TraceID().String(),
		PlatformName:       platformName,
		AppCluster:         appCluster,
		InstanceName:       instanceName,
		ParentSpanId:       span.ParentSpanID().String(),
		SpanId:             span.SpanID().String(),
		ServiceName:        serviceName,
		SpanName:           span.Name(),
		SpanKind:           span.Kind().String(),
		FaultKind:          faultKind,
		ResourceAttributes: attributesToMap(*attributes),
		SpanAttributes:     attributesToMap(span.Attributes()),
	}
	duration := span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).Nanoseconds()
	return &spanTreeItem{fault, duration, nil, 0}
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
