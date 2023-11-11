package spanaggregationexporter

import "github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/internal"

type Config struct {
	ClickHouseConfig internal.ClickHouseConfig `mapstructure:"clickhouse"`
	CacheConfig      internal.CacheConfig      `mapstructure:"cache"`
	BatchConfig      internal.BatchConfig      `mapstructure:"batch"`
}
