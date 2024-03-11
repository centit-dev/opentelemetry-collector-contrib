package internal

import (
	"context"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/systemparameter"
)

type SystemParameterRepository interface {
	FindByCode(ctx context.Context, code string) ([]string, error)
	Shutdown(ctx context.Context) error
}

type SystemParameterRepositoryImpl struct {
	client *DatabaseClient
}

func CreateSystemParameterRepository(client *DatabaseClient) *SystemParameterRepositoryImpl {
	return &SystemParameterRepositoryImpl{client}
}

func (repository *SystemParameterRepositoryImpl) FindByCode(ctx context.Context, code string) ([]string, error) {
	systemParameter, err := repository.client.delegate.SystemParameter.Query().
		Where(systemparameter.ID(code)).
		Select(systemparameter.FieldValue).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return systemParameter.Value, nil
}

func (repository *SystemParameterRepositoryImpl) Shutdown(ctx context.Context) error {
	return repository.client.delegate.Close()
}
