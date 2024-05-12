package nginxprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	conventions "go.opentelemetry.io/collector/semconv/v1.22.0"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

var np = CreateTraceProcessor()

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
