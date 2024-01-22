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
			values := make([]*ent.ApplicationStructure, len(items))
			for i, item := range items {
				values[i], ok = item.(*ent.ApplicationStructure)
				if !ok {
					return nil, nil
				}
			}
			err := service.repository.SaveAll(ctx, values)
			return values, err
		}).
		Retry(1, func(err error) bool { return true }).
		DoOnError(func(err error) {
			service.logger.Sugar().Errorf("cannot save application structures: %v", err)
		})

	service.refreshObservable = rxgo.FromChannel(service.refreshQueue)
	service.refreshObservable.
		BufferWithTimeOrCount(rxgo.WithDuration(time.Second*time.Duration(service.batchIntervalInSeconds)), service.batchSize)

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

func (service *ApplicationStructureService) deprecate(ctx context.Context) error {
	return service.repository.DeleteOutdated(ctx)
}

func (service *ApplicationStructureService) ConsumeAttributes(ctx context.Context, tuples map[string]*structureTuple) {
	// load cache
	if service.cache.Len() == 0 {
		items, err := service.repository.FindAll(ctx)
		if err != nil {
			return
		}
		for _, item := range items {
			service.cache.Add(item.ID, item)
		}
	}

	validDate := time.Now().AddDate(service.ttlInDays, 0, 0)
	for child, value := range tuples {
		item, ok := service.cache.Get(child)
		if ok && validDate.After(item.ValidDate) {
			service.refreshQueue <- rxgo.Of(item)
		} else {
			item = &ent.ApplicationStructure{
				ID:         child,
				ParentCode: value.parentCode,
				Level:      int(value.level),
				ValidDate:  validDate,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			service.cache.Add(child, item)
			service.batchQueue <- rxgo.Of(item)
		}
	}
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
