package internal

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	scopeNameKey         = "Scope[\"name\"]"
	scopeVersionKey      = "Scope[\"version\"]"
	statusCodeKey        = "StatusCode"
	spanSource           = "Span"
	metricSource         = "Metric"
	logSource            = "Log"
	queryValueTypeString = "S"
	queryValueTypeNumber = "N"
)

type dataPointSlice interface {
	Len() int
	At(i int) dataPoint
}

type dataPoint interface {
	Attributes() pcommon.Map
}

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
			k = fmt.Sprintf("ResourceAttributes[%s]", k)
			exporter.consumeAttribute(ctx, tuples, spanSource, k, v)
			return true
		})

		scopeSpans := resourceSpan.ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			scopeSpan := scopeSpans.At(j)
			spans := scopeSpan.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				exporter.consumeAttribute(ctx, tuples, spanSource, statusCodeKey, pcommon.NewValueStr(span.Status().Code().String()))
				spanAttributes := span.Attributes()
				spanAttributes.Range(func(k string, v pcommon.Value) bool {
					k = fmt.Sprintf("SpanAttributes[%s]", k)
					exporter.consumeAttribute(ctx, tuples, spanSource, k, v)
					return true
				})
			}
		}
	}
	exporter.service.ConsumeAttributes(ctx, tuples)
	return nil
}

func (exporter *MetadataExporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	tuples := make(map[string]*tuple)

	resourceMetrics := md.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		resourceMetric := resourceMetrics.At(i)
		resourceAttributes := resourceMetric.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("ResourceAttributes[%s]", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
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
					exporter.consumeGaugeDataPoints(ctx, tuples, metric.Gauge().DataPoints())
				case pmetric.MetricTypeSum:
					exporter.consumeGaugeDataPoints(ctx, tuples, metric.Sum().DataPoints())
				case pmetric.MetricTypeSummary:
					exporter.consumeSummaryDataPoints(ctx, tuples, metric.Summary().DataPoints())
				case pmetric.MetricTypeHistogram:
					exporter.consumeHistogramDataPoints(ctx, tuples, metric.Histogram().DataPoints())
				case pmetric.MetricTypeExponentialHistogram:
					exporter.consumeExpotentialHistogramDataPoints(ctx, tuples, metric.ExponentialHistogram().DataPoints())
				default:
					continue
				}
			}
		}
	}
	exporter.service.ConsumeAttributes(ctx, tuples)

	return nil
}

func (exporter *MetadataExporter) consumeGaugeDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.NumberDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes[%s]", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeSummaryDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.SummaryDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes[%s]", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeHistogramDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.HistogramDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes[%s]", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeExpotentialHistogramDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.ExponentialHistogramDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes[%s]", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	tuples := make(map[string]*tuple)

	resourceLogs := ld.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		resourceLog := resourceLogs.At(i)
		resourceAttributes := resourceLog.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Resource[%s]", k)
			exporter.consumeAttribute(ctx, tuples, logSource, k, v)
			return true
		})

		scopeLogs := resourceLog.ScopeLogs()
		for j := 0; j < scopeLogs.Len(); j++ {
			scopeLog := scopeLogs.At(j)
			logs := scopeLog.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				log := logs.At(k)
				log.Attributes().Range(func(k string, v pcommon.Value) bool {
					k = fmt.Sprintf("Attributes[%s]", k)
					exporter.consumeAttribute(ctx, tuples, logSource, k, v)
					return true
				})
			}
		}
	}
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
		name:         k,
		spansValid:   source == spanSource,
		metricsValid: source == metricSource,
		logsValid:    source == logSource,
		value:        value,
		valueType:    valueType,
	}
	_, ok := tuples[tuple.hash()]
	if !ok {
		tuples[tuple.hash()] = tuple
	}
}

func (exporter *MetadataExporter) Shutdown(ctx context.Context) error {
	return exporter.service.Shutdown(ctx)
}
