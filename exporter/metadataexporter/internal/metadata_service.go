package internal

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	"go.uber.org/zap"
)

type tuple struct {
	name         string
	spansValid   bool
	metricsValid bool
	logsValid    bool
	value        string
	valueType    string
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
	ConsumeAttribute(ctx context.Context, tuple tuple)
	Shutdown(ctx context.Context) error
}

type MetadataServiceImpl struct {
	cache                  *expirable.LRU[string, *ent.QueryKey]
	batchSize              int
	batchIntervalInSeconds int
	batchQueue             chan rxgo.Item
	batchObservable        rxgo.Observable
	refreshQueue           chan rxgo.Item
	refreshObservable      rxgo.Observable
	expirationObservable   rxgo.Observable
	logger                 *zap.Logger
	queryKeyTtlInDays      int
	queryKeyRepository     QueryKeyRepository
	queryValueRepository   QueryValueRepository
	systemParameterService SystemParameterService
}

func CreateMetadataService(
	cacheConfig *CacheConfig,
	batchConfig *BatchConfig,
	ttl int,
	logger *zap.Logger,
	queryKeyRepository QueryKeyRepository,
	queryValueRepository QueryValueRepository,
	systemParameterService SystemParameterService,
) *MetadataServiceImpl {
	cache := expirable.NewLRU[string, *ent.QueryKey](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	return &MetadataServiceImpl{
		cache:                  cache,
		batchSize:              batchConfig.BatchSize,
		batchIntervalInSeconds: batchConfig.IntervalInSeconds,
		logger:                 logger,
		queryKeyTtlInDays:      ttl,
		queryKeyRepository:     queryKeyRepository,
		queryValueRepository:   queryValueRepository,
		systemParameterService: systemParameterService,
	}
}

func (service *MetadataServiceImpl) Start(ctx context.Context) {
	service.systemParameterService.Start(ctx)

	service.batchQueue = make(chan rxgo.Item, 100)
	service.batchObservable = rxgo.FromChannel(service.batchQueue)
	service.batchObservable.
		BufferWithTimeOrCount(rxgo.WithDuration(time.Second*time.Duration(service.batchIntervalInSeconds)), service.batchSize).
		Map(func(ctx context.Context, batch interface{}) (interface{}, error) {
			items, ok := batch.([]interface{})
			if !ok {
				return nil, nil
			}
			seen := make(map[string]struct{}, len(items))
			tuples := make([]*tuple, 0, len(items))
			for _, item := range items {
				tuple, ok := item.(tuple)
				if !ok {
					continue
				}
				if _, ok := seen[tuple.hash()]; ok {
					continue
				}
				seen[tuple.hash()] = struct{}{}
				tuples = append(tuples, &tuple)
			}
			err := service.consumeAttributes(ctx, tuples)
			return items, err
		}).
		Retry(1, func(err error) bool { return true }).
		DoOnError(func(err error) {
			service.logger.Sugar().Errorf("cannot upsert query keys: %v", err)
		})

	service.refreshQueue = make(chan rxgo.Item, 100)
	service.refreshObservable = rxgo.FromChannel(service.refreshQueue)
	service.refreshObservable.
		BufferWithTimeOrCount(rxgo.WithDuration(time.Hour), service.batchSize).
		Map(func(ctx context.Context, batch interface{}) (interface{}, error) {
			items, ok := batch.([]interface{})
			if !ok {
				return nil, nil
			}

			err := service.refreshAll(ctx, items)
			return items, err
		}).
		Retry(1, func(err error) bool { return true }).
		DoOnError(func(err error) {
			service.logger.Sugar().Errorf("cannot refresh query values: %v", err)
		})

	service.expirationObservable = rxgo.Interval(rxgo.WithDuration(time.Hour))
	service.expirationObservable.
		Map(func(ctx context.Context, i interface{}) (interface{}, error) {
			err := service.deprecate(ctx)
			if err != nil {
				return nil, err
			}
			return i, nil
		}).
		// don't bother retrying deprecate
		DoOnError(func(err error) {
			service.logger.Sugar().Errorf("cannot deprecate keys: %v", err)
		})
}

func (service *MetadataServiceImpl) consumeAttributes(ctx context.Context, tuples []*tuple) error {
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

	upserts := make([]*ent.QueryKey, 0, len(tuples))
	validDate := time.Now().AddDate(0, 0, service.queryKeyTtlInDays)
	for _, tuple := range tuples {
		// find out new keys
		queryKey, ok := service.cache.Get(tuple.name)
		if ok {
			updated := false
			// if the key exists, update the valids
			if !queryKey.SpansValid && tuple.spansValid {
				queryKey.SpansValid = true
				updated = true
			}
			if !queryKey.MetricsValid && tuple.metricsValid {
				queryKey.MetricsValid = true
				updated = true
			}
			if !queryKey.LogsValid && tuple.logsValid {
				queryKey.LogsValid = true
				updated = true
			}
			if updated {
				queryKey.ValidDate = validDate
				upserts = append(upserts, queryKey)
			} else if validDate.After(queryKey.ValidDate.AddDate(0, 0, 1)) {
				// update the valid date
				queryKey.ValidDate = validDate
				service.refreshQueue <- rxgo.Of(queryKey)
			}
		} else {
			// if the key does not exist, create a new one
			queryKey = &ent.QueryKey{
				Name:         tuple.name,
				Type:         tuple.valueType,
				SpansValid:   tuple.spansValid,
				MetricsValid: tuple.metricsValid,
				LogsValid:    tuple.logsValid,
				ValidDate:    validDate,
			}
			service.cache.Add(tuple.name, queryKey)
		}

		// find out new values
		newQueryValue := &ent.QueryValue{Value: tuple.value, ValidDate: validDate}
		if existed, ok := service.cache.Get(valueHash(queryKey.Name, newQueryValue.Value)); ok {
			// if the value exists, update the valid date
			if validDate.After(existed.ValidDate.AddDate(0, 0, 1)) {
				existed.ValidDate = validDate
				service.refreshQueue <- rxgo.Of(existed)
			}
		} else {
			// if the value does not exist, create a new one
			queryKey.Edges.Values = append(queryKey.Edges.Values, newQueryValue)
			service.cache.Add(valueHash(queryKey.Name, newQueryValue.Value), queryKey)
			// it may look like we push multiple times for one key after seeing different values
			// but actually, all values are updated with the same key
			// so when upsertAll, updating one key will update all values
			upserts = append(upserts, queryKey)
		}
	}

	_, _, err = service.upsertAll(ctx, upserts)
	return err
}

func (service *MetadataServiceImpl) upsertAll(ctx context.Context, queryKeys []*ent.QueryKey) ([]*ent.QueryKey, []*ent.QueryKey, error) {
	if len(queryKeys) == 0 {
		return nil, nil, nil
	}
	toCreate := make([]*ent.QueryKey, 0, len(queryKeys))
	toUpdate := make([]*ent.QueryKey, 0, len(queryKeys))
	added := make(map[string]bool, len(queryKeys))
	for _, queryKey := range queryKeys {
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

func (service *MetadataServiceImpl) refreshAll(ctx context.Context, items []interface{}) error {
	if len(items) == 0 {
		return nil
	}

	validDate := time.UnixMilli(0)
	keys := make([]int64, 0, len(items))
	values := make([]int64, 0, len(items))
	for _, item := range items {
		if queryValue, ok := item.(*ent.QueryValue); ok {
			validDate = queryValue.ValidDate
			values = append(values, queryValue.ID)
		} else if queryKey, ok := item.(*ent.QueryKey); ok {
			validDate = queryKey.ValidDate
			keys = append(keys, queryKey.ID)
		}
	}

	err := service.queryKeyRepository.RefreshAllById(ctx, keys, validDate)
	if err != nil {
		return err
	}

	return service.queryValueRepository.RefreshAllById(ctx, values, validDate)
}

func (service *MetadataServiceImpl) deprecate(ctx context.Context) error {
	err := service.queryValueRepository.DeleteOutdated(ctx)
	if err != nil {
		return err
	}
	return service.queryKeyRepository.DeleteOutdated(ctx)
}

func (service *MetadataServiceImpl) ConsumeAttribute(ctx context.Context, tuple tuple) {
	if !service.systemParameterService.ShouldRecord(ctx, tuple.name) {
		return
	}
	service.batchQueue <- rxgo.Of(tuple)
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
	if service.refreshObservable != nil {
		_, cancel := service.refreshObservable.Connect(ctx)
		cancel()
	}
	if service.refreshQueue != nil {
		close(service.refreshQueue)
	}
	if service.expirationObservable != nil {
		_, cancel := service.expirationObservable.Connect(ctx)
		cancel()
	}
	return service.queryKeyRepository.Shutdown(ctx)
}
