package client

import (
	"context"
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/pkg/spangroup"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
	"go.uber.org/zap"
)

const (
	appTypeKey = "app.type"
)

type AppTypeService struct {
	logger *zap.Logger
	client DatabaseClient
	groups *spangroup.SpanGroups
	ticker *time.Ticker
}

func CreateAppTypeService(client DatabaseClient, cacheTtlMinutes time.Duration, logger *zap.Logger) *AppTypeService {
	ticker := time.NewTicker(cacheTtlMinutes * time.Minute)
	service := &AppTypeService{logger, client, nil, ticker}
	go func() {
		// build cache asynchronously for every 5 minutes so the first few batches won't be tagged and blocked
		defer ticker.Stop()
		for ; ; <-ticker.C {
			logger.Info("Building cache")
			service.buildCache(context.Background())
		}
	}()
	return service
}

func (service *AppTypeService) buildCache(context context.Context) {
	if service.client == nil {
		service.logger.Error("Error when building cache: database client is nil")
		return
	}
	records, err := service.client.FindAllDefinitions(context)
	if err != nil {
		service.logger.Sugar().Errorf("Error when querying groups: %s\n", err)
		return
	}
	data := make(map[*spangroup.SpanGroupDefinitions]string)
	for _, record := range records {
		definitions := spangroup.SpanGroupDefinitions{}
		for _, condition := range record.SpanConditions {
			definitions = append(definitions, spangroup.SpanGroupDefinition{
				Column: condition.Column,
				Op:     condition.Op,
				Value:  spangroup.CreateDefinitionValue(condition.Value),
			})
		}
		data[&definitions] = record.Name
	}
	service.groups = spangroup.CreateSpanGroup(data)
}

// implement processorhelper.ProcessTracesFunc
func (service *AppTypeService) ProcessTraces(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
	slice := traces.ResourceSpans()
	for i := 0; i < slice.Len(); i++ {
		resource := slice.At(i)
		resourceAttributes := resource.Resource().Attributes()

		batches := resource.ScopeSpans()
		for j := 0; j < batches.Len(); j++ {
			scopeSpans := batches.At(j)
			scope := scopeSpans.Scope()
			batch := scopeSpans.Spans()
			for k := 0; k < batch.Len(); k++ {
				span := batch.At(k)
				// TODO process span asynchronizely
				service.processSpan(&resourceAttributes, &scope, &span)
			}
		}
	}

	return traces, nil
}

func (service *AppTypeService) processSpan(resources *pcommon.Map, scope *pcommon.InstrumentationScope, span *ptrace.Span) {
	attributes := span.Attributes()
	if service.groups.IsEmpty() {
		return
	}
	// check db.system
	_, ok := attributes.Get(conventions.AttributeDBSystem)
	if !ok {
		return
	}
	queries := make(map[string]interface{}, resources.Len()+attributes.Len())
	// not all required
	resources.Range(func(k string, v pcommon.Value) bool {
		queries[k] = v.AsRaw()
		return true
	})
	// not all required
	attributes.Range(func(k string, v pcommon.Value) bool {
		queries[k] = v.AsRaw()
		return true
	})
	groups := service.groups.Get(&queries)
	if len(groups) > 0 {
		attributes.PutStr(appTypeKey, groups[0])
	}
}

// implement Shutdown from component.Component
func (service *AppTypeService) Shutdown(_ context.Context) error {
	service.ticker.Stop()
	return service.client.Shutdown()
}
