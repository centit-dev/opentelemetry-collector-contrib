package client

import (
	"testing"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// disable for integration test
func DisableTestCreateClient(t *testing.T) {
	config := &PostgresConfig{
		Host:     "host.docker.internal",
		Port:     25432,
		User:     "postgres",
		Pass:     "password",
		Database: "postgres",
	}
	logger, _ := zap.NewDevelopment()
	_, err := CreateClient(config, logger)
	if err != nil {
		t.Errorf("Error when creating client: %s\n", err)
	}
}
