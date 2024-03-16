package metadataexporter

import (
	"context"
	"fmt"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/internal"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	Type             = "metadata"
	TracesStability  = component.StabilityLevelAlpha
	MetricsStability = component.StabilityLevelAlpha
	LogsStability    = component.StabilityLevelAlpha
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		Type,
		createDefaultConfig,
		exporter.WithTraces(createTraceMetadataExporter, TracesStability),
		exporter.WithMetrics(createMetricMetadataExporter, MetricsStability),
		exporter.WithLogs(createLogMetadataExporter, LogsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		CacheConfig: internal.CacheConfig{
			MaxSize:         10000,
			ExpireInMinutes: 10,
		},
		BatchConfig: internal.BatchConfig{
			BatchSize:         1000,
			IntervalInSeconds: 1,
		},
		TtlInDays: 30,
	}
}

func createTraceMetadataExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	exporter, err := createExporter(set, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create metadata exporter: %w", err)
	}

	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		exporter.ConsumeTraces,
		exporterhelper.WithStart(exporter.Start),
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}

func createMetricMetadataExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Metrics, error) {
	exporter, err := createExporter(set, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create metadata exporter: %w", err)
	}

	return exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		exporter.ConsumeMetrics,
		exporterhelper.WithStart(exporter.Start),
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}

func createLogMetadataExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	exporter, err := createExporter(set, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create metadata exporter: %w", err)
	}

	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		exporter.ConsumeLogs,
		exporterhelper.WithStart(exporter.Start),
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}

func createExporter(set exporter.CreateSettings, cfg component.Config) (*internal.MetadataExporter, error) {
	c := cfg.(*Config)
	client, err := internal.CreateClient(&c.PostgresConfig, set.Logger)
	if err != nil {
		return nil, err
	}
	queryKeyRepository := internal.CreateQueryKeyRepository(client)
	queryValueRepository := internal.CreateQueryValueRepository(client)
	systemParameterRepository := internal.CreateSystemParameterRepository(client)
	systemParameterService := internal.CreateSystemParameterService(set.Logger, systemParameterRepository)
	metadataService := internal.CreateMetadataService(
		&c.CacheConfig,
		&c.BatchConfig,
		c.TtlInDays,
		set.Logger,
		queryKeyRepository,
		queryValueRepository,
		systemParameterService,
	)

	appStructureRepository := internal.CreateApplicationStructureRepository(client)
	appStructureService := internal.CreateApplicationStructureService(
		&c.CacheConfig,
		&c.BatchConfig,
		c.TtlInDays,
		set.Logger,
		appStructureRepository,
	)
	return internal.CreateMetadataExporter(metadataService, appStructureService), nil
}
