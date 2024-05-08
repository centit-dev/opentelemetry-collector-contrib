package logstashexporter

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/client/http"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/scheme"
)

type LogstashClient interface {
	Handle(ctx context.Context, msg *scheme.LogsOutput) (err error)
}

func InitClient(cfg *Config) (LogstashClient, error) {
	switch cfg.Type {
	case "http":
		return http.NewLsHttpClient(&cfg.HttpCfg), nil
	default:
		return nil, fmt.Errorf("invalid logstash client type")
	}
}
