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
	Type            = "metadata"
	TracesStability = component.StabilityLevelAlpha
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		Type,
		createDefaultConfig,
		exporter.WithTraces(createMetadataExporter, TracesStability),
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
		QueryKeyTtlInDays: 90,
	}
}

func createMetadataExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	c := cfg.(*Config)
	client, err := internal.CreateClient(&c.PostgresConfig, set.Logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create postgres sql client: %w", err)
	}
	queryKeyRepository := internal.CreateQueryKeyRepository(client)
	service := internal.CreateMetadataService(&c.CacheConfig, &c.BatchConfig, c.QueryKeyTtlInDays, set.Logger, queryKeyRepository)
	exporter := internal.CreateMetadataExporter(service)

	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		exporter.ConsumeTraces,
		exporterhelper.WithStart(exporter.Start),
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}
