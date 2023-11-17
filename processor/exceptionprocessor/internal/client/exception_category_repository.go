package client

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptioncategory"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptiondefinition"
)

type ExceptionCategoryRepository interface {
	FindAllDefinitions(context context.Context) ([]*ent.ExceptionDefinition, error)
	Shutdown(context context.Context) error
}

type ExceptionCategoryRepositoryImpl struct {
	client *PostgresClient
}

func NewExceptionCategoryRepository(client *PostgresClient) *ExceptionCategoryRepositoryImpl {
	return &ExceptionCategoryRepositoryImpl{client}
}

func (repo *ExceptionCategoryRepositoryImpl) FindAllDefinitions(context context.Context) ([]*ent.ExceptionDefinition, error) {
	return repo.client.delegate.ExceptionDefinition.Query().
		Where(
			exceptiondefinition.HasExceptionCategoryWith(exceptioncategory.IsValid(1)),
			exceptiondefinition.IsValid(1)).
		WithExceptionCategory().
		All(context)
}

func (repo *ExceptionCategoryRepositoryImpl) Shutdown(context context.Context) error {
	return repo.client.Shutdown()
}
