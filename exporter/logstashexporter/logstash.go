package logstashexporter

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/scheme"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/pdata/plog"
)

type logstashExporter struct {
	config *Config
	client LogstashClient
	logger *zap.Logger
}

func newLogstashExporter(cfg *Config, cli LogstashClient, logger *zap.Logger) *logstashExporter {
	return &logstashExporter{config: cfg, client: cli, logger: logger}
}

func (le *logstashExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) (err error) {
	// todo should it need to limit the sending count of output ?
	output := scheme.NewHttpOutputFromPlogs(ld)
	le.logger.Debug("send output to logstash",
		zap.Int("output size", len(output.Logs)))
	if !le.config.DryRun {
		err = le.client.Handle(ctx, output)
	}
	if err != nil {
		le.logger.Error("send output to logstash fail", zap.Error(err))
	}
	return
}
