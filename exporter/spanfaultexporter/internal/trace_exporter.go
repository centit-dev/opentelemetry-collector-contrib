package internal

import (
	"context"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
	"go.opentelemetry.io/collector/component"
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
	faults := make([]*ent.SpanFault, 0, traces.SpanCount())
	slice := traces.ResourceSpans()
	for i := 0; i < slice.Len(); i++ {
		resource := slice.At(i)
		resourceAttributes := resource.Resource().Attributes()
		var platformName string
		if platformNameValue, exists := resourceAttributes.Get(platFormNameKey); exists {
			platformName = platformNameValue.Str()
		}
		var clusterName string
		if clusterNameValue, exists := resourceAttributes.Get(conventions.AttributeK8SClusterName); exists {
			clusterName = clusterNameValue.Str()
		}
		var instanceName string
		if instanceNameValue, exists := resourceAttributes.Get(conventions.AttributeK8SNodeName); exists {
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
				fault := t.buildSpanFault(ctx, &span, platformName, clusterName, instanceName, serviceName)
				faults = append(faults, fault)
			}
		}
	}

	return t.service.Save(ctx, faults)
}

func (t *TraceExporter) buildSpanFault(ctx context.Context, span *ptrace.Span, platformName string, clusterName string, instanceName string, serviceName string) *ent.SpanFault {
	var faultKind string
	if faultKindValue, exists := span.Attributes().Get(faultKindKey); exists {
		faultKind = faultKindValue.Str()
	}
	return &ent.SpanFault{
		Timestamp:    span.StartTimestamp().AsTime(),
		TraceId:      span.TraceID().String(),
		PlatformName: platformName,
		ClusterName:  clusterName,
		InstanceName: instanceName,
		ParentSpanId: span.ParentSpanID().String(),
		ID:           span.SpanID().String(),
		ServiceName:  serviceName,
		SpanName:     span.Name(),
		FaultKind:    faultKind,
	}
}

func (t *TraceExporter) Shutdown(ctx context.Context) error {
	return t.service.Shutdown(ctx)
}
