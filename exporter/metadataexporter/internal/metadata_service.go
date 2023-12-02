package internal

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	"go.uber.org/zap"
)

const (
	queryKeySourceResource = "Resource"
	queryKeySourceSpan     = "Span"

	queryValueTypeString = "S"
	queryValueTypeNumber = "N"
)

type tuple struct {
	name      string
	source    string
	value     string
	valueType string
}

func (t *tuple) hash() string {
	return valueHash(t.name, t.value)
}

type CacheConfig struct {
	MaxSize         int `mapstructure:"max_size"`
	ExpireInMinutes int `mapstructure:"expire_in_minutes"`
}

type BatchConfig struct {
	BatchSize         int `mapstructure:"batch_size"`
	IntervalInSeconds int `mapstructure:"interval_in_seconds"`
}

type MetadataService interface {
	Start(ctx context.Context)
	ConsumeAttributes(ctx context.Context, tuples map[string]*tuple)
	Shutdown(ctx context.Context) error
}

type MetadataServiceImpl struct {
	cache                  *expirable.LRU[string, *ent.QueryKey]
	batchSize              int
	batchIntervalInSeconds int
	batchQueue             chan rxgo.Item
	batchObservable        rxgo.Observable
	logger                 *zap.Logger
	queryKeyTtlInDays      int
	queryKeyRepository     QueryKeyRepository
}

func CreateMetadataService(cacheConfig *CacheConfig, batchConfig *BatchConfig, ttl int, logger *zap.Logger, queryKeyRepository QueryKeyRepository) *MetadataServiceImpl {
	cache := expirable.NewLRU[string, *ent.QueryKey](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	return &MetadataServiceImpl{
		cache:                  cache,
		batchSize:              batchConfig.BatchSize,
		batchIntervalInSeconds: batchConfig.IntervalInSeconds,
		logger:                 logger,
		queryKeyTtlInDays:      ttl,
		queryKeyRepository:     queryKeyRepository,
	}
}

func (service *MetadataServiceImpl) Start(ctx context.Context) {
	service.batchQueue = make(chan rxgo.Item, 100)
	service.batchObservable = rxgo.FromChannel(service.batchQueue)
	service.batchObservable.
		BufferWithTimeOrCount(rxgo.WithDuration(time.Second*time.Duration(service.batchIntervalInSeconds)), service.batchSize).
		Map(func(ctx context.Context, batch interface{}) (interface{}, error) {
			items, ok := batch.([]interface{})
			if !ok {
				return nil, nil
			}
			_, _, err := service.upsertAll(ctx, items)
			return items, err
		}).
		Retry(1, func(err error) bool { return true }).
		DoOnError(func(err error) {
			service.logger.Sugar().Errorf("cannot upsert query keys: %v", err)
		})
}

func (service *MetadataServiceImpl) upsertAll(ctx context.Context, queryKeys []interface{}) ([]*ent.QueryKey, []*ent.QueryKey, error) {
	if len(queryKeys) == 0 {
		return nil, nil, nil
	}
	toCreate := make([]*ent.QueryKey, 0, len(queryKeys))
	toUpdate := make([]*ent.QueryKey, 0, len(queryKeys))
	added := make(map[string]bool, len(queryKeys))
	for _, item := range queryKeys {
		queryKey, ok := item.(*ent.QueryKey)
		if !ok {
			continue
		}
		if added[queryKey.Name] {
			continue
		}
		added[queryKey.Name] = true
		if queryKey.ID == 0 {
			toCreate = append(toCreate, queryKey)
		} else {
			toUpdate = append(toUpdate, queryKey)
		}
	}
	service.queryKeyRepository.RefreshState(ctx, toCreate)
	created, err := service.queryKeyRepository.CreateAll(ctx, toCreate)
	if err != nil {
		service.logger.Sugar().Errorf("cannot create query keys: %v", err)
		return nil, nil, err
	}
	service.queryKeyRepository.RefreshState(ctx, toUpdate)
	updated, err := service.queryKeyRepository.UpdateAll(ctx, toUpdate)
	if err != nil {
		service.logger.Sugar().Errorf("cannot update query keys: %v", err)
		return nil, nil, err
	}
	return created, updated, nil
}

func (service *MetadataServiceImpl) ConsumeAttributes(ctx context.Context, tuples map[string]*tuple) {
	// lazily build query key cache
	names := make([]string, 0, len(tuples))
	for _, tuple := range tuples {
		names = append(names, tuple.name)
	}
	queryKeys, err := service.queryKeyRepository.FindAllByNames(ctx, names)
	if err != nil {
		service.logger.Sugar().Errorf("cannot retrieve query keys by names: %v", err)
	} else if len(queryKeys) > 0 {
		for _, queryKey := range queryKeys {
			service.cache.Add(queryKey.Name, queryKey)
			for _, queryValue := range queryKey.Edges.Values {
				service.cache.Add(valueHash(queryKey.Name, queryValue.Value), queryKey)
			}
		}
	}

	// find out new keys and new values
	// if the key and the value are existing in the cache, update the valid date
	// if the key is existing in the cache, add the new value
	// if the key is not existing in the cache, create a new entry and add the new value
	for _, tuple := range tuples {
		queryKey, ok := service.cache.Get(tuple.name)
		validDate := time.Now().AddDate(0, 0, service.queryKeyTtlInDays)
		if ok {
			queryKey.ValidDate = validDate
		} else {
			queryKey = &ent.QueryKey{
				Name:      tuple.name,
				Type:      tuple.valueType,
				Source:    tuple.source,
				ValidDate: validDate,
			}
			service.cache.Add(tuple.name, queryKey)
		}

		newQueryValue := &ent.QueryValue{Value: tuple.value}
		_, ok = service.cache.Get(valueHash(queryKey.Name, newQueryValue.Value))
		if ok {
			continue
		}
		queryKey.Edges.Values = append(queryKey.Edges.Values, newQueryValue)
		service.cache.Add(valueHash(queryKey.Name, newQueryValue.Value), queryKey)
		service.batchQueue <- rxgo.Of(queryKey)
	}
}

func (service *MetadataServiceImpl) Shutdown(ctx context.Context) error {
	if service.cache != nil {
		service.cache.Purge()
	}
	if service.batchObservable != nil {
		_, cancel := service.batchObservable.Connect(ctx)
		cancel()
	}
	if service.batchQueue != nil {
		close(service.batchQueue)
	}
	return service.queryKeyRepository.Shutdown(ctx)
}
