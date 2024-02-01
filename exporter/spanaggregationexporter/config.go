package spanaggregationexporter

import (
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/internal"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
	configretry.BackOffConfig      `mapstructure:"retry_on_failure"`
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`

	ClickHouseConfig internal.ClickHouseConfig `mapstructure:"clickhouse"`
	CacheConfig      internal.CacheConfig      `mapstructure:"cache"`
	BatchConfig      internal.BatchConfig      `mapstructure:"batch"`
}
