package client

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"
	"go.uber.org/zap"
)

const faultKindKey = "fault.kind"

type TraceProcessor struct {
	logger  *zap.Logger
	service FaultKindService
}

func CreateTraceProcessor(logger *zap.Logger, service FaultKindService) *TraceProcessor {
	return &TraceProcessor{logger, service}
}

func (p *TraceProcessor) Start(ctx context.Context, _ component.Host) error {
	p.service.Start(ctx)
	return nil
}

func (p *TraceProcessor) ProcessTraces(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
	if !p.service.IsConfigured(ctx) {
		return traces, nil
	}
	slice := traces.ResourceSpans()
	for i := 0; i < slice.Len(); i++ {
		resourceSpans := slice.At(i)
		resourceAttributes := resourceSpans.Resource().Attributes().AsRaw()

		scopeSpans := resourceSpans.ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			spans := scopeSpans.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				p.processTrace(ctx, &resourceAttributes, &span)
			}
		}
	}
	return traces, nil
}

func (p *TraceProcessor) processTrace(ctx context.Context, resourceAttributes *map[string]interface{}, span *ptrace.Span) {
	spanAttributes := span.Attributes().AsRaw()
	faultKind := p.service.MatchFaultKind(ctx, resourceAttributes, &spanAttributes)
	if faultKind != "" {
		span.Attributes().PutStr(faultKindKey, faultKind)
	} else if p.hasException(span) {
		span.Attributes().PutStr(faultKindKey, SystemFault.String())
	}
}

func (TraceProcessor) hasException(span *ptrace.Span) bool {
	events := span.Events()
	for i := 0; i < events.Len(); i++ {
		event := events.At(i)
		attributes := event.Attributes()
		if _, ok := attributes.Get(conventions.AttributeExceptionType); ok {
			return true
		}
	}
	return false
}

func (p *TraceProcessor) Shutdown(ctx context.Context) error {
	return p.service.Shutdown(ctx)
}
