package exceptionprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/internal/client"
)

const (
	Type            = "exception"
	ServiceCacheKey = "service"
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		Type,
		createDefaultConfig,
		processor.WithTraces(createTraceProcessor, component.StabilityLevelAlpha),
		processor.WithLogs(createLogsProcessor, component.StabilityLevelAlpha),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		CacheTtlMinutes: 5,
	}
}

func createTraceProcessor(ctx context.Context, params processor.CreateSettings, cfg component.Config, nextConsumer consumer.Traces) (processor.Traces, error) {
	service, err := createService(cfg, params.Logger)
	if err != nil {
		return nil, err
	}
	// create trace processor with process func
	return processorhelper.NewTracesProcessor(ctx, params, cfg, nextConsumer, service.ProcessTraces, processorhelper.WithShutdown(service.Shutdown))
}

func createLogsProcessor(ctx context.Context, params processor.CreateSettings, cfg component.Config, nextConsumer consumer.Logs) (processor.Logs, error) {
	service, err := createService(cfg, params.Logger)
	if err != nil {
		return nil, err
	}
	// create logs processor with process func
	return processorhelper.NewLogsProcessor(ctx, params, cfg, nextConsumer, service.ProcessLogs, processorhelper.WithShutdown(service.Shutdown))
}

func createService(cfg component.Config, logger *zap.Logger) (*client.ExceptionCategoryService, error) {
	config := cfg.(*Config)
	databaseClient, err := client.CreateClient(&config.Postgres, logger)
	if err != nil {
		logger.Sugar().Errorf("Error when creating %s: %s\n", Type, err)
		return nil, err
	}
	repo := client.NewExceptionCategoryRepository(databaseClient)

	// create category service
	service := client.CreateCategoryService(repo, config.CacheTtlMinutes, logger)
	return service, nil
}
