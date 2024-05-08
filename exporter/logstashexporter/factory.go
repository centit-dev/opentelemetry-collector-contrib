package logstashexporter

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

var (
	Type = component.MustNewType("logstash")
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		Type,
		createDefaultConfig,
		exporter.WithLogs(newLogstashExporterfunc, component.StabilityLevelBeta))
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func newLogstashExporterfunc(ctx context.Context, params exporter.CreateSettings, config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validate fail, err: %v", err)
	}
	client, err := InitClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("init client fail, err: %v", err)
	}
	debugExporter := logstashExporter{config: cfg, client: client, logger: params.Logger}
	return exporterhelper.NewLogsExporter(ctx, params, config,
		debugExporter.ConsumeLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// todo configurable?
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 3 * time.Second}),
	)
}
