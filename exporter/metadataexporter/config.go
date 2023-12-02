package metadataexporter

import "github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/internal"

type Config struct {
	PostgresConfig internal.PostgresConfig `mapstructure:"postgres"`
	CacheConfig    internal.CacheConfig    `mapstructure:"cache"`
	BatchConfig    internal.BatchConfig    `mapstructure:"batch"`
	// TTL in days for query keys and values before they are removed. Default is 90 days.
	QueryKeyTtlInDays int `mapstructure:"ttl_in_days"`
}
