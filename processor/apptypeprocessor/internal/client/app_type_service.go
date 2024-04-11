package client

import (
	"context"
	"fmt"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/pkg/spangroup"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.22.0"
	"go.uber.org/zap"
)

type softwareType int16

const (
	serverSoftware = "server.software"
	serverUrl      = "server.url"
	appType        = "application.type"
)

const (
	typeApplication softwareType = iota + 1
	typeServerSoftware
)

type AppTypeService struct {
	logger *zap.Logger
	client DatabaseClient
	groups *spangroup.SpanGroups
	ticker *time.Ticker

	records map[string]*ent.SoftwareDefinition
}

func CreateAppTypeService(client DatabaseClient, cacheTtlMinutes time.Duration, logger *zap.Logger) *AppTypeService {
	ticker := time.NewTicker(cacheTtlMinutes * time.Minute)
	service := &AppTypeService{logger, client, nil, ticker, make(map[string]*ent.SoftwareDefinition)}
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
		groupName := fmt.Sprintf("%s.%d", record.Name, record.Type)
		data[&definitions] = groupName
		service.records[groupName] = record
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
			batch := scopeSpans.Spans()
			for k := 0; k < batch.Len(); k++ {
				span := batch.At(k)
				// TODO process span asynchronizely
				service.processSpan(&resourceAttributes, &span)
			}
		}
	}

	return traces, nil
}

func (service *AppTypeService) processSpan(resources *pcommon.Map, span *ptrace.Span) {
	if service.groups.IsEmpty() {
		return
	}

	attributes := span.Attributes()
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
	service.setAttributes(&attributes, groups)
}

func (service *AppTypeService) setAttributes(attributes *pcommon.Map, groups []string) {
	for _, group := range groups {
		record, ok := service.records[group]
		if !ok {
			continue
		}
		switch record.Type {
		case int16(typeApplication):
			attributes.PutStr(appType, record.Name)
		case int16(typeServerSoftware):
			attributes.PutStr(serverSoftware, record.Name)
			if url, ok := attributes.Get(conventions.AttributeDBConnectionString); ok {
				attributes.PutStr(serverUrl, url.Str())
			}
		}
	}
}

// implement Shutdown from component.Component
func (service *AppTypeService) Shutdown(_ context.Context) error {
	service.ticker.Stop()
	return service.client.Shutdown()
}
