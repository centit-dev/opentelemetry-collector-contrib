package internal

import (
	"context"
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/queryvalue"
)

type QueryValueRepository interface {
	RefreshAllById(ctx context.Context, ids []int64, validDate time.Time) error
	DeleteOutdated(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type QueryValueRepositoryImpl struct {
	client *DatabaseClient
}

func CreateQueryValueRepository(client *DatabaseClient) *QueryValueRepositoryImpl {
	return &QueryValueRepositoryImpl{client}
}

func (repository *QueryValueRepositoryImpl) RefreshAllById(ctx context.Context, ids []int64, validDate time.Time) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := repository.client.delegate.QueryValue.Update().
		Where(queryvalue.IDIn(ids...)).
		SetValidDate(validDate).
		Save(ctx)
	return err
}

func (repository *QueryValueRepositoryImpl) DeleteOutdated(ctx context.Context) error {
	_, err := repository.client.delegate.QueryValue.Delete().
		Where(queryvalue.ValidDateLT(time.Now())).
		Exec(ctx)
	return err
}

func (repository *QueryValueRepositoryImpl) Shutdown(ctx context.Context) error {
	return repository.client.delegate.Close()
}
