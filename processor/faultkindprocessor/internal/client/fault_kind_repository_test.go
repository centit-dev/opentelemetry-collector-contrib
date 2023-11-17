package client

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

// disabled for integration test
func DisabledTestFaultKindRepositoryImpl_FindFaultKindDefinition(t *testing.T) {
	config := &PostgresConfig{Host: "host.docker.internal", Port: 25432, User: "postgres", Pass: "password"}
	client, err := CreateClient(config, zap.NewExample())
	if err != nil {
		t.Errorf("error when building postgres client: %v", err)
	}
	repository := CreateFaultKindRepositoryImpl(client)

	definitions, err := repository.FindFaultKindDefinitions(context.Background())
	if err != nil && definitions.System == nil && definitions.Business == nil {
		t.Errorf("err when accessing fault kind definitions: %v %v", err, definitions)
	}
}
