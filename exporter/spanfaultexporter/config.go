package spanfaultexporter

import "github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/internal"

type Config struct {
	ClickHouseConfig internal.ClickHouseConfig `mapstructure:"clickhouse"`
	BatchConfig      internal.BatchConfig      `mapstructure:"batch"`
}
