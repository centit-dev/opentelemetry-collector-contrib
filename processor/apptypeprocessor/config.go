package apptypeprocessor

import (
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/internal/client"
)

type Config struct {
	Postgres        client.PostgresConfig `mapstructure:"postgres"`
	CacheTtlMinutes time.Duration         `mapstructure:"cache_ttl_minutes"`
}
