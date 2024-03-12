package internal

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	"go.uber.org/zap"
)

type applicationStructureLevel int

const (
	levelPlatform           applicationStructureLevel = 1
	levelApplicationCluster applicationStructureLevel = 2
	levelInstance           applicationStructureLevel = 3
)

type structureTuple struct {
	parentCode string
	code       string
	level      applicationStructureLevel
}

type ApplicationStructureService struct {
	cache                  *expirable.LRU[string, *ent.ApplicationStructure]
	batchSize              int
	batchIntervalInSeconds int
	batchQueue             chan rxgo.Item
	batchObservable        rxgo.Observable
	refreshQueue           chan rxgo.Item
	refreshObservable      rxgo.Observable
	expirationObservable   rxgo.Observable
	logger                 *zap.Logger
	ttlInDays              int
	repository             ApplicationStructureRepository
}

func CreateApplicationStructureService(
	cacheConfig *CacheConfig,
	batchConfig *BatchConfig,
	ttl int,
	logger *zap.Logger,
	repository ApplicationStructureRepository,
) *ApplicationStructureService {
	cache := expirable.NewLRU[string, *ent.ApplicationStructure](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	return &ApplicationStructureService{
		cache:                  cache,
		batchSize:              batchConfig.BatchSize,
		batchIntervalInSeconds: batchConfig.IntervalInSeconds,
		batchQueue:             make(chan rxgo.Item, batchConfig.BatchSize),
		refreshQueue:           make(chan rxgo.Item, batchConfig.BatchSize),
		logger:                 logger,
		ttlInDays:              ttl,
		repository:             repository,
	}
}

func (service *ApplicationStructureService) Start(ctx context.Context) {
	service.batchObservable = rxgo.FromChannel(service.batchQueue)
	service.batchObservable.
		BufferWithTimeOrCount(rxgo.WithDuration(time.Second*time.Duration(service.batchIntervalInSeconds)), service.batchSize).
		Map(func(ctx context.Context, batch interface{}) (interface{}, error) {
			items, ok := batch.([]interface{})
			if !ok {
				return nil, nil
			}
			seen := make(map[string]struct{})
			values := make([]*structureTuple, 0, len(items))
			for _, item := range items {
				value, ok := item.(structureTuple)
				if !ok {
					continue
				}
				if _, ok := seen[value.code]; ok {
					continue
				}
				seen[value.code] = struct{}{}
				values = append(values, &value)
			}
			err := service.upserts(ctx, values)
			return values, err
		}).
		Retry(1, func(err error) bool { return true }).
		DoOnError(func(err error) {
			service.logger.Sugar().Errorf("cannot save application structures: %v", err)
		})

	service.refreshObservable = rxgo.FromChannel(service.refreshQueue)
	service.refreshObservable.
		BufferWithTimeOrCount(rxgo.WithDuration(time.Second*time.Duration(service.batchIntervalInSeconds)), service.batchSize).
		Map(func(ctx context.Context, items interface{}) (interface{}, error) {
			return items, nil
		}).
		// don't bother retrying deprecate
		DoOnError(func(err error) {
			service.logger.Sugar().Errorf("cannot deprecate keys: %v", err)
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

func (service *ApplicationStructureService) upserts(ctx context.Context, tuples []*structureTuple) error {
	// load cache
	if service.cache.Len() == 0 {
		items, err := service.repository.FindAll(ctx)
		if err != nil {
			return err
		}
		for _, item := range items {
			service.cache.Add(item.ID, item)
		}
	}

	validDate := time.Now().AddDate(service.ttlInDays, 0, 0)
	upserts := make([]*ent.ApplicationStructure, 0, len(tuples))
	for _, tuple := range tuples {
		item, ok := service.cache.Get(tuple.code)
		if ok && validDate.After(item.ValidDate) {
			service.refreshQueue <- rxgo.Of(item)
		} else {
			item = &ent.ApplicationStructure{
				ID:         tuple.code,
				ParentCode: tuple.parentCode,
				Level:      int(tuple.level),
				ValidDate:  validDate,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			service.cache.Add(tuple.code, item)
			upserts = append(upserts, item)
		}
	}

	return service.repository.SaveAll(ctx, upserts)
}

func (service *ApplicationStructureService) deprecate(ctx context.Context) error {
	return service.repository.DeleteOutdated(ctx)
}

func (service *ApplicationStructureService) ConsumeAttribute(ctx context.Context, tuple structureTuple) {
	service.batchQueue <- rxgo.Of(tuple)
}

func (service *ApplicationStructureService) Shutdown(ctx context.Context) error {
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
	return service.repository.Shutdown(ctx)
}
