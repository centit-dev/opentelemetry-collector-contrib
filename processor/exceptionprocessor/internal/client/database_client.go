package client

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptioncategory"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptiondefinition"
	"go.uber.org/zap"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Pass     string `mapstructure:"pass"`
	Database string `mapstructure:"database"`
}

type DatabaseClient interface {
	FindAllDefinitions(context context.Context) ([]*ent.ExceptionDefinition, error)
	Shutdown() error
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
		logger.Sugar().Errorf("failed opening connection to postgres: %v", err)
		return &PostgresClient{}, err
	}
	logger.Sugar().Infof("connected to postgres: %v:%v", config.Host, config.Port)
	return &PostgresClient{delegate}, nil
}

func (client *PostgresClient) FindAllDefinitions(context context.Context) ([]*ent.ExceptionDefinition, error) {
	return client.delegate.ExceptionDefinition.Query().
		Where(
			exceptiondefinition.HasExceptionCategoryWith(exceptioncategory.IsValid(1)),
			exceptiondefinition.IsValid(1)).
		WithExceptionCategory().
		All(context)
}

func (client *PostgresClient) Shutdown() error {
	// TODO: close is idempotent?
	return client.delegate.Close()
}
