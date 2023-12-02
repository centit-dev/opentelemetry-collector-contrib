package internal

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	"go.uber.org/zap"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Pass     string `mapstructure:"pass"`
	Database string `mapstructure:"database"`
	Debug    bool   `mapstructure:"debug"`
}

type DatabaseClient struct {
	delegate *ent.Client
}

func CreateClient(config *PostgresConfig, logger *zap.Logger) (*DatabaseClient, error) {
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Pass, config.Host, config.Port, config.Database)
	delegate, err := ent.Open("postgres", url)
	if err != nil {
		logger.Sugar().Errorf("failed opening connection to postgres: %v", err)
		return &DatabaseClient{}, err
	}
	logger.Sugar().Infof("connected to postgres: %v:%v", config.Host, config.Port)
	if config.Debug {
		delegate = delegate.Debug()
	}
	return &DatabaseClient{delegate}, nil
}
