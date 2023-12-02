package internal

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type MetadataExporter struct {
	service MetadataService
}

func CreateMetadataExporter(service MetadataService) *MetadataExporter {
	return &MetadataExporter{service}
}

func (exporter *MetadataExporter) Start(ctx context.Context, _ component.Host) error {
	exporter.service.Start(ctx)
	return nil
}

func (exporter *MetadataExporter) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	tuples := make(map[string]*tuple)

	resourceSpans := td.ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		resourceSpan := resourceSpans.At(i)
		resourceAttributes := resourceSpan.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			exporter.consumeAttribute(ctx, tuples, queryKeySourceResource, k, v)
			return true
		})

		scopeSpans := resourceSpan.ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			scopeSpan := scopeSpans.At(j)
			spans := scopeSpan.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				spanAttributes := span.Attributes()
				spanAttributes.Range(func(k string, v pcommon.Value) bool {
					exporter.consumeAttribute(ctx, tuples, queryKeySourceSpan, k, v)
					return true
				})
			}
		}
	}
	exporter.service.ConsumeAttributes(ctx, tuples)
	return nil
}

func (exporter *MetadataExporter) consumeAttribute(ctx context.Context, tuples map[string]*tuple, source string, k string, v pcommon.Value) {
	valueType := queryValueTypeNumber
	value := fmt.Sprint(v.Int())
	if v.Type() == pcommon.ValueTypeStr {
		valueType = queryValueTypeString
		value = v.Str()
	}
	tuple := &tuple{
		name:      k,
		source:    source,
		value:     value,
		valueType: valueType,
	}
	_, ok := tuples[tuple.hash()]
	if !ok {
		tuples[tuple.hash()] = tuple
	}
}

func (exporter *MetadataExporter) Shutdown(ctx context.Context) error {
	return exporter.service.Shutdown(ctx)
}
