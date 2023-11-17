package client

import (
	"context"

	"github.com/teanoon/opentelemetry-collector-contrib/processor/faultkindprocessor/ent/schema"
	"github.com/teanoon/opentelemetry-collector-contrib/processor/faultkindprocessor/ent/systemparameter"
)

const faulkKindCode = "error_rules"

type FaultKindRepository interface {
	FindFaultKindDefinitions(ctx context.Context) (*schema.FaultKindDefinitions, error)
	Shutdown(ctx context.Context) error
}

type FaultKindRepositoryImpl struct {
	client *PostgresClient
}

func CreateFaultKindRepositoryImpl(client *PostgresClient) *FaultKindRepositoryImpl {
	return &FaultKindRepositoryImpl{client}
}

func (r *FaultKindRepositoryImpl) FindFaultKindDefinitions(ctx context.Context) (*schema.FaultKindDefinitions, error) {
	params, err := r.client.delegate.SystemParameter.Query().
		Where(systemparameter.IDEQ(faulkKindCode)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return params.Value, nil
}

func (r *FaultKindRepositoryImpl) Shutdown(ctx context.Context) error {
	return r.client.Shutdown(ctx)
}
