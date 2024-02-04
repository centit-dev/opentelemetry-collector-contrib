package internal

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"
)

const (
	scopeNameKey         = "Scope[\"name\"]"
	scopeVersionKey      = "Scope[\"version\"]"
	spanNameKey          = "SpanName"
	statusCodeKey        = "StatusCode"
	spanSource           = "Span"
	metricSource         = "Metric"
	logSource            = "Log"
	queryValueTypeString = "S"
	queryValueTypeNumber = "N"

	attributeServicePlatform = "service.platform"
)

type MetadataExporter struct {
	metadataService     MetadataService
	appStructureService *ApplicationStructureService
}

func CreateMetadataExporter(service MetadataService, appStructureService *ApplicationStructureService) *MetadataExporter {
	return &MetadataExporter{service, appStructureService}
}

func (exporter *MetadataExporter) Start(ctx context.Context, _ component.Host) error {
	exporter.metadataService.Start(ctx)
	exporter.appStructureService.Start(ctx)
	return nil
}

func (exporter *MetadataExporter) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	resourceSpans := td.ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		resourceSpan := resourceSpans.At(i)
		resourceAttributes := resourceSpan.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("ResourceAttributes['%s']", k)
			exporter.consumeAttribute(ctx, spanSource, k, v)
			return true
		})

		platform, ok := resourceAttributes.Get(attributeServicePlatform)
		if ok {
			exporter.consumeStructureAttributes(ctx, platform.Str(), "", levelPlatform)

			appCluster, ok := resourceAttributes.Get(conventions.AttributeK8SDeploymentName)
			if ok {
				exporter.consumeStructureAttributes(ctx, appCluster.Str(), platform.Str(), levelApplicationCluster)

				podName, ok := resourceAttributes.Get(conventions.AttributeK8SPodName)
				if ok {
					exporter.consumeStructureAttributes(ctx, podName.Str(), appCluster.Str(), levelInstance)
				}
			}
		}

		scopeSpans := resourceSpan.ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			scopeSpan := scopeSpans.At(j)
			spans := scopeSpan.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				exporter.consumeAttribute(ctx, spanSource, spanNameKey, pcommon.NewValueStr(span.Name()))
				exporter.consumeAttribute(ctx, spanSource, statusCodeKey, pcommon.NewValueStr(span.Status().Code().String()))
				spanAttributes := span.Attributes()
				spanAttributes.Range(func(k string, v pcommon.Value) bool {
					k = fmt.Sprintf("SpanAttributes['%s']", k)
					exporter.consumeAttribute(ctx, spanSource, k, v)
					return true
				})
			}
		}
	}
	return nil
}

func (exporter *MetadataExporter) consumeStructureAttributes(ctx context.Context, code string, parentCode string, level applicationStructureLevel) {
	item := structureTuple{
		parentCode: parentCode,
		code:       code,
		level:      level,
	}
	exporter.appStructureService.ConsumeAttribute(ctx, item)
}

func (exporter *MetadataExporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	resourceMetrics := md.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		resourceMetric := resourceMetrics.At(i)
		resourceAttributes := resourceMetric.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("ResourceAttributes['%s']", k)
			exporter.consumeAttribute(ctx, metricSource, k, v)
			return true
		})

		scopeMetrics := resourceMetric.ScopeMetrics()
		for j := 0; j < scopeMetrics.Len(); j++ {
			scopeMetric := scopeMetrics.At(j)
			metrics := scopeMetric.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)
				switch metric.Type() {
				case pmetric.MetricTypeGauge:
					exporter.consumeGaugeDataPoints(ctx, metric.Gauge().DataPoints())
				case pmetric.MetricTypeSum:
					exporter.consumeGaugeDataPoints(ctx, metric.Sum().DataPoints())
				case pmetric.MetricTypeSummary:
					exporter.consumeSummaryDataPoints(ctx, metric.Summary().DataPoints())
				case pmetric.MetricTypeHistogram:
					exporter.consumeHistogramDataPoints(ctx, metric.Histogram().DataPoints())
				case pmetric.MetricTypeExponentialHistogram:
					exporter.consumeExpotentialHistogramDataPoints(ctx, metric.ExponentialHistogram().DataPoints())
				default:
					continue
				}
			}
		}
	}

	return nil
}

func (exporter *MetadataExporter) consumeGaugeDataPoints(ctx context.Context, dps pmetric.NumberDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
			exporter.consumeAttribute(ctx, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeSummaryDataPoints(ctx context.Context, dps pmetric.SummaryDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
			exporter.consumeAttribute(ctx, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeHistogramDataPoints(ctx context.Context, dps pmetric.HistogramDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
			exporter.consumeAttribute(ctx, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeExpotentialHistogramDataPoints(ctx context.Context, dps pmetric.ExponentialHistogramDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
			exporter.consumeAttribute(ctx, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	resourceLogs := ld.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		resourceLog := resourceLogs.At(i)
		resourceAttributes := resourceLog.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Resource['%s']", k)
			exporter.consumeAttribute(ctx, logSource, k, v)
			return true
		})

		scopeLogs := resourceLog.ScopeLogs()
		for j := 0; j < scopeLogs.Len(); j++ {
			scopeLog := scopeLogs.At(j)
			logs := scopeLog.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				log := logs.At(k)
				log.Attributes().Range(func(k string, v pcommon.Value) bool {
					k = fmt.Sprintf("Attributes['%s']", k)
					exporter.consumeAttribute(ctx, logSource, k, v)
					return true
				})
			}
		}
	}

	return nil
}

func (exporter *MetadataExporter) consumeAttribute(ctx context.Context, source string, k string, v pcommon.Value) {
	valueType := queryValueTypeNumber
	value := fmt.Sprint(v.Int())
	if v.Type() == pcommon.ValueTypeStr {
		valueType = queryValueTypeString
		value = v.Str()
	}
	tuple := &tuple{
		name:         k,
		spansValid:   source == spanSource,
		metricsValid: source == metricSource,
		logsValid:    source == logSource,
		value:        value,
		valueType:    valueType,
	}
	exporter.metadataService.ConsumeAttribute(ctx, *tuple)
}

func (exporter *MetadataExporter) Shutdown(ctx context.Context) error {
	err1 := exporter.metadataService.Shutdown(ctx)
	err2 := exporter.appStructureService.Shutdown(ctx)
	return mergeErrors(err1, err2)
}

func mergeErrors(err1, err2 error) error {
	if err1 != nil && err2 != nil {
		return fmt.Errorf("%s; %s", err1, err2)
	} else if err1 != nil {
		return err1
	} else {
		return err2
	}
}
