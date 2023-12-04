package internal

import (
	"context"
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/querykey"
)

type QueryKeyRepository interface {
	FindAll(ctx context.Context) ([]*ent.QueryKey, error)
	FindAllByNames(ctx context.Context, names []string) ([]*ent.QueryKey, error)
	RefreshAllById(ctx context.Context, ids []int64, validDate time.Time) error
	DeleteOutdated(ctx context.Context) error
	// the collector may execute in parallel, so the key/value cache may be out of date
	// double check the database state before creating duplicates
	RefreshState(ctx context.Context, querKeys []*ent.QueryKey) error
	// create keys and values in batches
	CreateAll(ctx context.Context, queryKeys []*ent.QueryKey) ([]*ent.QueryKey, error)
	// update keys and values in batches
	UpdateAll(ctx context.Context, queryKeys []*ent.QueryKey) ([]*ent.QueryKey, error)
	Shutdown(ctx context.Context) error
}

type QueryKeyRepositoryImpl struct {
	client *DatabaseClient
}

func CreateQueryKeyRepository(client *DatabaseClient) *QueryKeyRepositoryImpl {
	return &QueryKeyRepositoryImpl{client}
}

func (repository *QueryKeyRepositoryImpl) FindAll(ctx context.Context) ([]*ent.QueryKey, error) {
	return repository.client.delegate.QueryKey.Query().
		Where(querykey.ValidDateGTE(time.Now())).
		WithValues().
		All(ctx)
}

func (repository *QueryKeyRepositoryImpl) FindAllByNames(ctx context.Context, names []string) ([]*ent.QueryKey, error) {
	if len(names) == 0 {
		return nil, nil
	}

	return repository.client.delegate.QueryKey.Query().
		Where(querykey.ValidDateGTE(time.Now()), querykey.NameIn(names...)).
		WithValues().
		All(ctx)
}

func (repository *QueryKeyRepositoryImpl) RefreshAllById(ctx context.Context, ids []int64, validDate time.Time) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := repository.client.delegate.QueryKey.Update().
		Where(querykey.IDIn(ids...)).
		SetValidDate(validDate).
		Save(ctx)
	return err
}

func (repository *QueryKeyRepositoryImpl) DeleteOutdated(ctx context.Context) error {
	_, err := repository.client.delegate.QueryKey.Delete().
		Where(querykey.ValidDateLT(time.Now())).
		Exec(ctx)
	return err
}

func (repository *QueryKeyRepositoryImpl) RefreshState(ctx context.Context, querKeys []*ent.QueryKey) error {
	if len(querKeys) == 0 {
		return nil
	}

	// find existing keys and values
	pendingKeys := make([]string, 0, len(querKeys))
	for _, queryKey := range querKeys {
		pendingKeys = append(pendingKeys, queryKey.Name)
	}
	existingQueryKeys, err := repository.FindAllByNames(ctx, pendingKeys)
	if err != nil {
		return err
	}
	valueCount := 0
	for _, queryKey := range existingQueryKeys {
		valueCount += len(queryKey.Edges.Values)
	}
	existingKeyMap := make(map[string]*ent.QueryKey, len(existingQueryKeys))
	for _, queryKey := range existingQueryKeys {
		existingKeyMap[queryKey.Name] = queryKey
	}
	existingValueMap := make(map[string]*ent.QueryValue, valueCount)
	for _, queryKey := range existingQueryKeys {
		for _, queryValue := range queryKey.Edges.Values {
			existingValueMap[valueHash(queryKey.Name, queryValue.Value)] = queryValue
		}
	}

	// update keys and values cache
	for _, queryKey := range querKeys {
		if existing, ok := existingKeyMap[queryKey.Name]; ok {
			queryKey.ID = existing.ID
			queryKey.Name = existing.Name
			queryKey.Type = existing.Type
			queryKey.Source = existing.Source
			queryKey.ValidDate = existing.ValidDate
			queryKey.CreateTime = existing.CreateTime
			queryKey.UpdateTime = existing.UpdateTime
		}
		for _, queryValue := range queryKey.Edges.Values {
			if existing, ok := existingValueMap[valueHash(queryKey.Name, queryValue.Value)]; ok {
				queryValue.ID = existing.ID
				queryValue.KeyID = existing.KeyID
				queryValue.Value = existing.Value
				queryValue.ValidDate = existing.ValidDate
				queryValue.CreateTime = existing.CreateTime
				queryValue.UpdateTime = existing.UpdateTime
				queryValue.Edges.Key = existing.Edges.Key
			}
		}
	}

	return nil
}

func (repository *QueryKeyRepositoryImpl) CreateAll(ctx context.Context, queryKeys []*ent.QueryKey) ([]*ent.QueryKey, error) {
	if len(queryKeys) == 0 {
		return nil, nil
	}

	tx, err := repository.client.delegate.Tx(ctx)
	if err != nil {
		return nil, err
	}

	// create keys first
	toCreateKeys := make([]*ent.QueryKeyCreate, 0, len(queryKeys))
	for _, queryKey := range queryKeys {
		if queryKey.ID != 0 {
			continue
		}
		queryKey.CreateTime = time.Now()
		queryKey.UpdateTime = time.Now()
		toCreateKeys = append(toCreateKeys, tx.QueryKey.Create().
			SetName(queryKey.Name).
			SetType(queryKey.Type).
			SetSource(queryKey.Source).
			SetValidDate(queryKey.ValidDate).
			SetCreateTime(queryKey.CreateTime).
			SetUpdateTime(queryKey.UpdateTime))
	}

	// all keys exist, it becomes a updateAll
	if len(toCreateKeys) == 0 {
		return repository.UpdateAll(ctx, queryKeys)
	}

	createdKeys, err := tx.QueryKey.CreateBulk(toCreateKeys...).Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// update keys' ids
	createdKeyMap := make(map[string]*ent.QueryKey, len(createdKeys))
	for _, queryKey := range createdKeys {
		createdKeyMap[queryKey.Name] = queryKey
	}
	for _, queryKey := range queryKeys {
		createdKey, ok := createdKeyMap[queryKey.Name]
		if !ok {
			continue
		}
		queryKey.ID = createdKey.ID
		queryKey.ValidDate = createdKey.ValidDate
	}

	err = repository.createQueryValues(ctx, tx, queryKeys)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return queryKeys, nil
}

func (repository *QueryKeyRepositoryImpl) UpdateAll(ctx context.Context, queryKeys []*ent.QueryKey) ([]*ent.QueryKey, error) {
	if len(queryKeys) == 0 {
		return nil, nil
	}

	tx, err := repository.client.delegate.Tx(ctx)
	if err != nil {
		return nil, err
	}

	// update keys with new valid date and update time
	toUpdateKeys := make([]int64, 0, len(queryKeys))
	newValidDate := time.UnixMilli(0)
	for _, queryKey := range queryKeys {
		if queryKey.ID == 0 {
			continue
		}
		// all keys should have the same valid date
		newValidDate = queryKey.ValidDate
		toUpdateKeys = append(toUpdateKeys, queryKey.ID)
	}
	err = tx.QueryKey.Update().Where(querykey.IDIn(toUpdateKeys...)).
		SetValidDate(newValidDate).
		SetUpdateTime(time.Now()).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// add new values for the old keys
	err = repository.createQueryValues(ctx, tx, queryKeys)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// commit
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return queryKeys, nil
}

func (repository *QueryKeyRepositoryImpl) createQueryValues(ctx context.Context, tx *ent.Tx, queryKeys []*ent.QueryKey) error {
	if len(queryKeys) == 0 {
		return nil
	}

	// build query value creates
	valueCount := 0
	for _, queryKey := range queryKeys {
		valueCount += len(queryKey.Edges.Values)
	}
	queryValues := make([]*ent.QueryValueCreate, 0, valueCount)
	for _, queryKey := range queryKeys {
		if queryKey.ID == 0 {
			continue
		}
		for _, queryValue := range queryKey.Edges.Values {
			if queryValue.ID != 0 {
				continue
			}
			queryValue.CreateTime = time.Now()
			queryValue.UpdateTime = time.Now()
			queryValue.KeyID = queryKey.ID
			queryValue.Edges.Key = queryKey
			queryValues = append(queryValues, tx.QueryValue.Create().
				SetValue(queryValue.Value).
				SetValidDate(queryValue.ValidDate).
				SetCreateTime(queryValue.CreateTime).
				SetUpdateTime(queryValue.UpdateTime).
				SetKeyID(queryKey.ID))
		}
	}

	// insert values
	createdValues, err := tx.QueryValue.CreateBulk(queryValues...).Save(ctx)
	if err != nil {
		return err
	}

	// update values' ids
	createdValueMap := make(map[string]*ent.QueryValue, len(createdValues))
	for _, queryValue := range createdValues {
		createdValueMap[queryValue.Value] = queryValue
	}
	for _, queryKey := range queryKeys {
		for _, queryValue := range queryKey.Edges.Values {
			createdValue, ok := createdValueMap[queryValue.Value]
			if !ok {
				continue
			}
			queryValue.ID = createdValue.ID
		}
	}

	return nil
}

func (repository *QueryKeyRepositoryImpl) Shutdown(_ context.Context) error {
	return repository.client.delegate.Close()
}
