package client

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
	"go.uber.org/zap"
)

const (
	eventAttributeExceptionName         = "exception"
	spanAttributeSpanNameKey            = "span.name"
	spanAttributeExceptionDefinitionKey = "exception.definition.id"
)

type ExceptionCategoryService struct {
	logger       *zap.Logger
	repository   ExceptionCategoryRepository
	exceptionIds map[string]string
	ticker       *time.Ticker
}

func CreateCategoryService(repository ExceptionCategoryRepository, cacheTtlMinutes time.Duration, logger *zap.Logger) *ExceptionCategoryService {
	ticker := time.NewTicker(cacheTtlMinutes * time.Minute)
	service := &ExceptionCategoryService{logger, repository, make(map[string]string), ticker}
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

func (service *ExceptionCategoryService) buildCache(context context.Context) {
	if service.repository == nil {
		service.logger.Error("Error when building cache: database client is nil")
		return
	}
	records, err := service.repository.FindAllDefinitions(context)
	if err != nil {
		service.logger.Sugar().Errorf("Error when querying categories: %s\n", err)
		return
	}
	for _, definition := range records {
		service.exceptionIds[definition.LongName] = fmt.Sprint(definition.ID)
	}
}

// implement processorhelper.ProcessTracesFunc
func (service *ExceptionCategoryService) ProcessTraces(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
	if len(service.exceptionIds) == 0 {
		return traces, nil
	}

	slice := traces.ResourceSpans()
	for i := 0; i < slice.Len(); i++ {
		resource := slice.At(i)
		resourceAttributes := resource.Resource().Attributes()

		batches := resource.ScopeSpans()
		for j := 0; j < batches.Len(); j++ {
			batch := batches.At(j).Spans()
			for k := 0; k < batch.Len(); k++ {
				span := batch.At(k)
				// TODO process span asynchronizely
				service.processSpan(&resourceAttributes, &span)
			}
		}
	}

	return traces, nil
}

// implement processorhelper.ProcessLogsFunc
func (service *ExceptionCategoryService) ProcessLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
	if len(service.exceptionIds) == 0 {
		return logs, nil
	}

	slice := logs.ResourceLogs()
	for i := 0; i < slice.Len(); i++ {
		resource := slice.At(i)
		resourceAttributes := resource.Resource().Attributes()

		batches := resource.ScopeLogs()
		for j := 0; j < batches.Len(); j++ {
			batch := batches.At(j).LogRecords()
			for k := 0; k < batch.Len(); k++ {
				log := batch.At(k)
				// TODO process log asynchronizely
				service.processLog(&resourceAttributes, &log)
			}
		}
	}

	return logs, nil
}

// check if the span contains exception
// if yes, find the category and add exception.name and exception.type to the span
// if not, return the span unchanged
func (service *ExceptionCategoryService) processSpan(resources *pcommon.Map, span *ptrace.Span) {
	exceptionFullName := service.extractException(span)
	if exceptionFullName == "" {
		return
	}
	if id, ok := service.exceptionIds[exceptionFullName]; ok {
		span.Attributes().PutStr(spanAttributeExceptionDefinitionKey, id)
	}
}

func (service *ExceptionCategoryService) extractException(span *ptrace.Span) string {
	events := span.Events()
	for i := 0; i < events.Len(); i++ {
		event := events.At(i)
		if event.Name() == eventAttributeExceptionName {
			attributes := event.Attributes()
			if value, ok := attributes.Get(conventions.AttributeExceptionType); ok {
				return value.AsString()
			}
		}
	}
	return ""
}

// TODO logs doesn't have span name
func (service *ExceptionCategoryService) processLog(resources *pcommon.Map, log *plog.LogRecord) {
	if log.SeverityNumber() < plog.SeverityNumberError {
		return
	}
	attributes := log.Attributes()
	exceptionFullName, ok := attributes.Get(conventions.AttributeExceptionType)
	if !ok {
		return
	}
	exceptionFullNameValue := exceptionFullName.AsString()
	if id, ok := service.exceptionIds[exceptionFullNameValue]; ok {
		log.Attributes().PutStr(spanAttributeExceptionDefinitionKey, id)
	}
}

// implement Shutdown from component.Component
func (service *ExceptionCategoryService) Shutdown(ctx context.Context) error {
	service.ticker.Stop()
	return service.repository.Shutdown(ctx)
}
