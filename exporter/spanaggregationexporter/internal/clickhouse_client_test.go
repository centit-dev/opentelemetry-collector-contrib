package internal

import (
	"testing"

	"go.uber.org/zap"
)

func createRealClient() *ClickHouseClient {
	c, _ := CreateClient(&ClickHouseConfig{
		Endpoint: "host.docker.internal:29000",
		Database: "otel",
	}, zap.NewExample())
	return c
}

// disabled for integration test
func DisabledTestCreateClient(t *testing.T) {
	client := createRealClient()
	t.Run("CreateClientSuccessful", func(t *testing.T) {
		if client == nil {
			t.Errorf("CreateClient() = %v, want %v", client, nil)
		}
	})
}

// disabled for integration test
func DisabledTestClickHouseClient_Shutdown(t *testing.T) {
	client := createRealClient()
	t.Run("Shutdown twice", func(t *testing.T) {
		if err := client.Shutdown(); err != nil {
			t.Errorf("ClickHouseClient.Shutdown() error = %v", err)
		}
		if err := client.Shutdown(); err != nil {
			t.Errorf("ClickHouseClient.Shutdown() error = %v", err)
		}
	})
}
