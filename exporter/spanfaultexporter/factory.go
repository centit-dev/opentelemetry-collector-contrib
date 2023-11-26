package spanfaultexporter

import (
	"context"
	"fmt"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/internal"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	Type            = "spanfault"
	TracesStability = component.StabilityLevelAlpha
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		Type,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, TracesStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ClickHouseConfig: internal.ClickHouseConfig{
			DialTimeoutInSeconds:   30,
			MaxOpenConns:           10,
			MaxIdleConns:           5,
			ConnMaxLifetimeInHours: 1,
			Database:               "otel",
			Debug:                  false,
		},
		CacheConfig: internal.CacheConfig{MaxSize: 60_000 * 60, ExpireInMinutes: 5},
		BatchConfig: internal.BatchConfig{BatchSize: 1000, IntervalInMilliseconds: 1000},
	}
}

func createTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	c := cfg.(*Config)
	client, err := internal.CreateClient(&c.ClickHouseConfig, set.Logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create clickhouse sql client: %w", err)
	}
	repository := internal.CreateSpanFaultRepository(client, set.Logger)
	service := internal.CreateSpanFaultService(&c.CacheConfig, &c.BatchConfig, repository, set.Logger)
	exporter := internal.CreateTraceExporter(service, set.Logger)

	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		exporter.PushTraceData,
		exporterhelper.WithStart(exporter.Start),
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}
