package client

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/teanoon/opentelemetry-collector-contrib/processor/faultkindprocessor/ent"
	"go.uber.org/zap"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Pass     string `mapstructure:"pass"`
	Database string `mapstructure:"database"`
}

type PostgresClient struct {
	delegate *ent.Client
}

func CreateClient(config *PostgresConfig, logger *zap.Logger) (*PostgresClient, error) {
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Pass, config.Host, config.Port, config.Database)
	delegate, err := ent.Open(dialect.Postgres, url)
	if err != nil {
		return nil, err
	}
	logger.Sugar().Infof("connected to postgres: %v:%v", config.Host, config.Port)
	return &PostgresClient{delegate}, nil
}

func (c *PostgresClient) Shutdown(_ context.Context) error {
	return c.delegate.Close()
}
