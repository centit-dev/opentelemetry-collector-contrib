package processor

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"strings"
)

var _ sdktrace.SpanProcessor = (*SpanProcessor)(nil)

// SpanProcessor It tracks all active spans, and stores samples of spans based on latency for non errored spans,
// and samples for errored spans.
type SpanProcessor struct {
	handleDelayHistogram metric.Float64Histogram
}

// NewSpanProcessor returns a new SpanProcessor.
func NewSpanProcessor(m metric.Float64Histogram) *SpanProcessor {
	return &SpanProcessor{handleDelayHistogram: m}
}

// OnStart adds span as active and reports it.
func (ssm *SpanProcessor) OnStart(_ context.Context, span sdktrace.ReadWriteSpan) {
}

// OnEnd processes all spans and reports them.
func (ssm *SpanProcessor) OnEnd(span sdktrace.ReadOnlySpan) {
	code := span.Status().Code

	var status = "success"
	if code == codes.Error {
		//ss.errors.add(span)
		//return
		status = "fail"
	}

	latency := span.EndTime().Sub(span.StartTime())
	// In case of time skew or wrong time, sample as 0 latency.
	if latency < 0 {
		latency = 0
	}

	component_type, component_name := parseSpanName(span.Name())
	ssm.handleDelayHistogram.Record(context.Background(), latency.Seconds(), metric.WithAttributes(
		attribute.String("type", component_type),
		attribute.String("name", component_name),
		attribute.String("status", status),
	))
}

// Shutdown does nothing.
func (ssm *SpanProcessor) Shutdown(context.Context) error {
	// Do nothing
	return nil
}

// ForceFlush does nothing.
func (ssm *SpanProcessor) ForceFlush(context.Context) error {
	// Do nothing
	return nil
}

func parseSpanName(name string) (ctype string, cname string) {
	info := strings.Split(name, "/")
	if len(info) > 0 {
		ctype = info[0]
	}
	if len(info) > 1 {
		cname = info[1]
	}
	return
}
