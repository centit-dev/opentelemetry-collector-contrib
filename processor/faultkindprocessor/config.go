package faultkindprocessor

import (
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/processor/faultkindprocessor/internal/client"
)

type Config struct {
	Postgres        client.PostgresConfig `mapstructure:"postgres"`
	CacheTtlMinutes time.Duration         `mapstructure:"cache_ttl_minutes"`
}
