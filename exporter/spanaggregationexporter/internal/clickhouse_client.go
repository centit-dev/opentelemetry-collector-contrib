package internal

import (
	"time"

	"entgo.io/ent/dialect/sql"
	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"go.uber.org/zap"
)

type ClickHouseConfig struct {
	// Endpoint is the clickhouse endpoint.
	Endpoint string `mapstructure:"endpoint"`
	// DialTimeoutInSeconds is the dial timeout in seconds.
	DialTimeoutInSeconds int `mapstructure:"dial_timeout_in_seconds"`
	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// MaxIdleConns is the maximum number of connections in the idle connection pool.
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// ConnMaxLifetimeInHours is the maximum amount of time a connection may be reused.
	ConnMaxLifetimeInHours int `mapstructure:"conn_max_lifetime_in_hours"`
	// Database is the database name to export.
	Database string `mapstructure:"database"`
	// Username is the username to connect to clickhouse.
	Username string `mapstructure:"username"`
	// Password is the password to connect to clickhouse.
	Password string `mapstructure:"password"`
	// enable debug mode to print SQLs
	Debug bool `mapstructure:"debug"`
}

type ClickHouseClient struct {
	delegate *ent.Client
	driver   *sql.Driver
}

func CreateClient(config *ClickHouseConfig, logger *zap.Logger) (*ClickHouseClient, error) {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{config.Endpoint},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.Username,
			Password: config.Password,
		},
		DialTimeout: time.Second * time.Duration(config.DialTimeoutInSeconds),
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
			Level:  3,
		},
		Debug: config.Debug,
	})
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetimeInHours) * time.Hour)
	if err := db.Ping(); err != nil {
		logger.Sugar().Errorf("failed pinging clickhouse: %v", err)
		return nil, err
	}
	driver := sql.OpenDB("clickhouse", db)
	delegate := ent.NewClient(ent.Driver(driver))
	logger.Sugar().Infof("connected to clickhouse: %v %v", config.Endpoint, config.Database)
	if config.Debug {
		delegate = delegate.Debug()
	}
	return &ClickHouseClient{delegate, driver}, nil
}

func (client *ClickHouseClient) Shutdown() error {
	// TODO: close is idempotent?
	return client.delegate.Close()
}
