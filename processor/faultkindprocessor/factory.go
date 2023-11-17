package faultkindprocessor

import (
	"context"

	"github.com/teanoon/opentelemetry-collector-contrib/processor/faultkindprocessor/internal/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	Type = "faultkind"
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

	repository := client.CreateFaultKindRepositoryImpl(databaseClient)
	service := client.CreateFaultKindServiceImpl(params.Logger, repository, config.CacheTtlMinutes)
	processor := client.CreateTraceProcessor(service)
	// create trace processor with process func
	return processorhelper.NewTracesProcessor(
		ctx,
		params,
		cfg,
		nextConsumer,
		processor.ProcessTraces,
		processorhelper.WithStart(processor.Start),
		processorhelper.WithShutdown(processor.Shutdown))
}
