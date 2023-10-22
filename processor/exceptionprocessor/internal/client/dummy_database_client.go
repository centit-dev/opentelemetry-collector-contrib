package client

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// for testing purpose
type DummyDatabaseClient struct {
}

func (client *DummyDatabaseClient) FindAllDefinitions(context context.Context) ([]*ent.ExceptionDefinition, error) {
	return []*ent.ExceptionDefinition{}, nil
}

func (client *DummyDatabaseClient) Shutdown() error {
	return nil
}

func CreateTestService(cfg component.Config) (*ExceptionCategoryService, error) {
	logger, _ := zap.NewDevelopment()
	service := CreateCategoryService(&DummyDatabaseClient{}, 10, logger)
	return service, nil
}
