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
	tuples := make(map[string]*tuple)
	structureTuples := make(map[string]*structureTuple)

	resourceSpans := td.ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		resourceSpan := resourceSpans.At(i)
		resourceAttributes := resourceSpan.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("ResourceAttributes['%s']", k)
			exporter.consumeAttribute(ctx, tuples, spanSource, k, v)
			return true
		})

		platform, ok := resourceAttributes.Get(attributeServicePlatform)
		if ok {
			exporter.consumeStructureAttributes(ctx, structureTuples, platform.Str(), "", levelPlatform)

			appCluster, ok := resourceAttributes.Get(conventions.AttributeK8SDeploymentName)
			if ok {
				exporter.consumeStructureAttributes(ctx, structureTuples, appCluster.Str(), platform.Str(), levelApplicationCluster)

				podName, ok := resourceAttributes.Get(conventions.AttributeK8SPodName)
				if ok {
					exporter.consumeStructureAttributes(ctx, structureTuples, podName.Str(), appCluster.Str(), levelInstance)
				}
			}
		}

		scopeSpans := resourceSpan.ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			scopeSpan := scopeSpans.At(j)
			spans := scopeSpan.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				exporter.consumeAttribute(ctx, tuples, spanSource, spanNameKey, pcommon.NewValueStr(span.Name()))
				exporter.consumeAttribute(ctx, tuples, spanSource, statusCodeKey, pcommon.NewValueStr(span.Status().Code().String()))
				spanAttributes := span.Attributes()
				spanAttributes.Range(func(k string, v pcommon.Value) bool {
					k = fmt.Sprintf("SpanAttributes['%s']", k)
					exporter.consumeAttribute(ctx, tuples, spanSource, k, v)
					return true
				})
			}
		}
	}
	go func() {
		ctx := context.Background()
		exporter.metadataService.ConsumeAttributes(ctx, tuples)
		exporter.appStructureService.ConsumeAttributes(ctx, structureTuples)
	}()
	return nil
}

func (exporter *MetadataExporter) consumeStructureAttributes(ctx context.Context, tuples map[string]*structureTuple, code string, parentCode string, level applicationStructureLevel) {
	if _, ok := tuples[code]; ok {
		return
	}
	item := &structureTuple{
		parentCode: parentCode,
		level:      level,
	}
	tuples[code] = item
}

func (exporter *MetadataExporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	tuples := make(map[string]*tuple)

	resourceMetrics := md.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		resourceMetric := resourceMetrics.At(i)
		resourceAttributes := resourceMetric.Resource().Attributes()
		resourceAttributes.Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("ResourceAttributes['%s']", k)
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
	exporter.metadataService.ConsumeAttributes(ctx, tuples)

	return nil
}

func (exporter *MetadataExporter) consumeGaugeDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.NumberDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeSummaryDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.SummaryDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeHistogramDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.HistogramDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
			exporter.consumeAttribute(ctx, tuples, metricSource, k, v)
			return true
		})
	}
}

func (exporter *MetadataExporter) consumeExpotentialHistogramDataPoints(ctx context.Context, tuples map[string]*tuple, dps pmetric.ExponentialHistogramDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			k = fmt.Sprintf("Attributes['%s']", k)
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
			k = fmt.Sprintf("Resource['%s']", k)
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
					k = fmt.Sprintf("Attributes['%s']", k)
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
