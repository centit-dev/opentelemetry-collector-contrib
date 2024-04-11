package apptypeprocessor

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/internal/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

var (
	Type = component.MustNewType("apptype")
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		Type,
		createDefaultConfig,
		processor.WithTraces(createTraceProcessor, component.StabilityLevelAlpha),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		CacheTtlMinutes: 5,
	}
}

func createTraceProcessor(ctx context.Context, params processor.CreateSettings, cfg component.Config, nextConsumer consumer.Traces) (processor.Traces, error) {
	config := cfg.(*Config)
	databaseClient, err := client.CreateClient(&config.Postgres, params.Logger)
	if err != nil {
		params.Logger.Sugar().Errorf("Error when creating %s: %s\n", Type, err)
		return nil, err
	}

	service := client.CreateAppTypeService(databaseClient, config.CacheTtlMinutes, params.Logger)
	// create trace processor with process func
	return processorhelper.NewTracesProcessor(ctx, params, cfg, nextConsumer, service.ProcessTraces, processorhelper.WithShutdown(service.Shutdown))
}
