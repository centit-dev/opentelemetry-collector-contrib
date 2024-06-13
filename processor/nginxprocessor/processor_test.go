package nginxprocessor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	conventions "go.opentelemetry.io/collector/semconv/v1.22.0"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

var np = CreateTraceProcessor(nil)

func TestProcessTracesWithNginxSpansBaggage(t *testing.T) {
	// given a trace with nginx span
	traces := ptrace.NewTraces()
	resourceSpans := traces.ResourceSpans().AppendEmpty()
	scoped := resourceSpans.ScopeSpans().AppendEmpty()
	scoped.Scope().SetName(scopeNginx)
	span := scoped.Spans().AppendEmpty()
	span.Attributes().PutStr(conventions.AttributeK8SPodName, "nice-pod")

	expectBizCode := "13770321976"
	expectTransID := "98ASDF78SDF"
	baggage := fmt.Sprintf("biz_code=%s, trans_id=%s, good-key=good-value1;good-value2, more-key=more-value", expectBizCode, expectTransID)
	span.Attributes().PutStr(baggageKey, baggage)

	// when processResourceSpans is called
	traces, _ = np.ProcessTraces(context.Background(), traces)

	// then the span is collected in a pod name group
	assert.Equal(t, 1, traces.ResourceSpans().Len())
	processed := traces.ResourceSpans().At(0).Resource()
	podName, _ := processed.Attributes().Get(conventions.AttributeK8SPodName)
	assert.Equal(t, "nice-pod", podName.Str())
	assert.Equal(t, 1, traces.ResourceSpans().At(0).ScopeSpans().Len())
	processedScopeSpans := traces.ResourceSpans().At(0).ScopeSpans().At(0)
	assert.Equal(t, 1, processedScopeSpans.Spans().Len())
	processedSpan := processedScopeSpans.Spans().At(0)
	bizCode, ok := processedSpan.Attributes().Get(BizCodeKey)
	assert.Equal(t, true, ok)
	assert.Equal(t, expectBizCode, bizCode.AsString())
	transID, ok := processedSpan.Attributes().Get(TransIDKey)
	assert.Equal(t, true, ok)
	assert.Equal(t, expectTransID, transID.AsString())
}

func TestProcessTracesWithNginxSpans(t *testing.T) {
	// given a trace with nginx span
	traces := ptrace.NewTraces()
	resourceSpans := traces.ResourceSpans().AppendEmpty()
	scoped := resourceSpans.ScopeSpans().AppendEmpty()
	scoped.Scope().SetName(scopeNginx)
	span := scoped.Spans().AppendEmpty()
	span.Attributes().PutStr(conventions.AttributeK8SPodName, "nice-pod")

	// when processResourceSpans is called
	traces, _ = np.ProcessTraces(context.Background(), traces)

	// then the span is collected in a pod name group
	assert.Equal(t, 1, traces.ResourceSpans().Len())
	processed := traces.ResourceSpans().At(0).Resource()
	podName, _ := processed.Attributes().Get(conventions.AttributeK8SPodName)
	assert.Equal(t, "nice-pod", podName.Str())
}

func TestProcessTracesWithNonNginxSpans(t *testing.T) {
	// given a trace with non-nginx span
	traces := ptrace.NewTraces()
	resourceSpans := traces.ResourceSpans().AppendEmpty()
	scoped := resourceSpans.ScopeSpans().AppendEmpty()
	scoped.Scope().SetName("other")
	span := scoped.Spans().AppendEmpty()
	span.Attributes().PutStr(conventions.AttributeK8SPodName, "nice-pod")

	// when processResourceSpans is called
	traces, _ = np.ProcessTraces(context.Background(), traces)

	// then the span is not collected in a pod name group
	assert.Equal(t, 1, traces.ResourceSpans().Len())
	processed := traces.ResourceSpans().At(0).Resource()
	_, ok := processed.Attributes().Get(conventions.AttributeK8SPodName)
	assert.False(t, ok)
}

func TestProcessTracesWithMixedSpans(t *testing.T) {
	// given a trace with multiple nginx spans
	traces := ptrace.NewTraces()
	resourceSpans := traces.ResourceSpans().AppendEmpty()
	scoped := resourceSpans.ScopeSpans().AppendEmpty()
	scoped.Scope().SetName(scopeNginx)
	span := scoped.Spans().AppendEmpty()
	span.Attributes().PutStr(conventions.AttributeK8SPodName, "nice-pod")
	span = scoped.Spans().AppendEmpty()
	span.Attributes().PutStr(conventions.AttributeK8SPodName, "nice-pod")

	// and a scope with a non-nginx span
	scoped = resourceSpans.ScopeSpans().AppendEmpty()
	scoped.Scope().SetName("other")
	span = scoped.Spans().AppendEmpty()
	span.Attributes().PutStr(conventions.AttributeK8SPodName, "nice-pod2")

	// when processResourceSpans is called
	traces, _ = np.ProcessTraces(context.Background(), traces)

	// then the spans are collected in a pod name group
	assert.Equal(t, 2, traces.ResourceSpans().Len())
	resourceSpans0 := traces.ResourceSpans().At(0)
	resourceSpans1 := traces.ResourceSpans().At(1)
	var processed ptrace.ResourceSpans
	if _, ok := resourceSpans0.Resource().Attributes().Get(conventions.AttributeK8SPodName); ok {
		processed = resourceSpans0
	} else {
		processed = resourceSpans1
	}

	assert.Equal(t, 1, processed.ScopeSpans().Len())
	podName, _ := processed.Resource().Attributes().Get(conventions.AttributeK8SPodName)
	assert.Equal(t, "nice-pod", podName.Str())
}
