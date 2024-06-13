package nginxprocessor

import (
	"context"
	"go.uber.org/zap"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.22.0"
	"go.opentelemetry.io/otel/baggage"
)

const (
	scopeNginx      = "nginx"
	platformNameKey = "service.platform"

	baggageKey        = "source.baggage"
	baggageBizCodeKey = "biz_code"
	baggageTransIDKey = "trans_id"
	BizCodeKey        = "source.biz_code"
	TransIDKey        = "source.trans_id"
)

type scopeGroup struct {
	resourceSpans *ptrace.ResourceSpans
	scopeSpans    *ptrace.ScopeSpans
}

type nginxProcessor struct {
	logger *zap.Logger
}

func CreateTraceProcessor(logger *zap.Logger) *nginxProcessor {
	return &nginxProcessor{logger: logger}
}

func (np *nginxProcessor) ProcessTraces(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
	newTraces := ptrace.NewTraces()
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		resourceSpans := traces.ResourceSpans().At(i)
		np.processResourceSpans(&newTraces, &resourceSpans)
	}
	return newTraces, nil
}

// loop through all resource spans
// if a resource have a span from the nginx, collect the span in a pod name group.
// if a resource have no span from the nginx, keep it as it is.
// after the loop, create new resource spans with:
// 1. the resource with no span from the nginx
// 2. the resource with spans having the common pod name
func (np *nginxProcessor) processResourceSpans(newTraces *ptrace.Traces, resourceSpans *ptrace.ResourceSpans) {
	groups := make(map[string][]scopeGroup)

	origin := ptrace.NewResourceSpans()
	resourceSpans.Resource().CopyTo(origin.Resource())
	origin.SetSchemaUrl(resourceSpans.SchemaUrl())

	scopes := resourceSpans.ScopeSpans()
	for i := 0; i < scopes.Len(); i++ {
		scope := scopes.At(i)
		if podName, ok := np.isNginxScope(&scope); ok {
			// add the scope to the group by the pod name
			group, ok := groups[podName]
			if !ok {
				group = make([]scopeGroup, 0)
			}
			group = append(group, scopeGroup{
				resourceSpans: resourceSpans,
				scopeSpans:    &scope,
			})
			groups[podName] = group
		} else {
			newScope := origin.ScopeSpans().AppendEmpty()
			scope.CopyTo(newScope)
		}
	}

	if origin.ScopeSpans().Len() > 0 {
		reserved := newTraces.ResourceSpans().AppendEmpty()
		origin.CopyTo(reserved)
	}

	if len(groups) == 0 {
		return
	}

	// create new resource
	resourceSpansGroup := make(map[string]ptrace.ResourceSpans)
	for podName, group := range groups {
		created, ok := resourceSpansGroup[podName]
		if !ok {
			created = newTraces.ResourceSpans().AppendEmpty()
			resourceSpans.Resource().CopyTo(created.Resource())
			attributes := created.Resource().Attributes()
			rewriteServiceName(&attributes)
			attributes.PutStr(conventions.AttributeK8SPodName, podName)
			created.SetSchemaUrl(resourceSpans.SchemaUrl())
		}

		for _, g := range group {
			newScope := created.ScopeSpans().AppendEmpty()
			g.scopeSpans.CopyTo(newScope)
		}
	}
}

func (np *nginxProcessor) isNginxScope(scope *ptrace.ScopeSpans) (string, bool) {
	// only process nginx scope
	// other scopes are returned with its original resource
	if scope.Scope().Name() != scopeNginx {
		return "", false
	}

	spans := scope.Spans()
	// skip empty scope
	if spans.Len() == 0 {
		return "", false
	}

	for i := 0; i < spans.Len(); i++ {
		np.ParseBaggage(spans.At(i))
		if podName, ok := spans.At(i).Attributes().Get(conventions.AttributeK8SPodName); ok {
			return podName.Str(), true
		}
	}
	// skip if the pod name is not available
	return "", false
}

func (np *nginxProcessor) ParseBaggage(span ptrace.Span) {
	raw, ok := span.Attributes().Get(baggageKey)
	if !ok {
		return
	}
	value, err := baggage.Parse(raw.AsString())
	if err != nil {
		np.logger.Error("fail to parse baggage", zap.Error(err), zap.String("raw", raw.AsString()))
		return
	}
	bizCode := value.Member(baggageBizCodeKey)
	if bizCode.Value() != "" {
		span.Attributes().PutStr(BizCodeKey, bizCode.Value())
	}
	transID := value.Member(baggageTransIDKey)
	if transID.Value() != "" {
		span.Attributes().PutStr(TransIDKey, transID.Value())
	}
}

func rewriteServiceName(attributes *pcommon.Map) {
	serviceNameValue, ok := attributes.Get(conventions.AttributeServiceName)
	if !ok {
		return
	}
	serviceName := serviceNameValue.Str()
	if !strings.Contains(serviceName, ":") {
		return
	}
	platformServiceNames := strings.Split(serviceName, ":")
	if len(platformServiceNames) != 2 {
		return
	}
	platformName, serviceName := platformServiceNames[0], platformServiceNames[1]
	attributes.PutStr(conventions.AttributeServiceName, serviceName)
	attributes.PutStr(platformNameKey, platformName)
}
