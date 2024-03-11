package internal

import (
	"context"
	"time"

	"github.com/reactivex/rxgo/v2"
)

const metadataKeyBlacklist = "metadata_key_blacklist"

type SystemParameterService interface {
	Start(ctx context.Context)
	ShouldRecord(ctx context.Context, code string) bool
	Shutdown(ctx context.Context) error
}

type SystemParameterServiceImpl struct {
	cache             map[string]struct{}
	repository        SystemParameterRepository
	refreshObservable rxgo.Observable
}

func CreateSystemParameterService(repository SystemParameterRepository) *SystemParameterServiceImpl {
	return &SystemParameterServiceImpl{
		cache:      make(map[string]struct{}),
		repository: repository,
	}
}

func (service *SystemParameterServiceImpl) Start(ctx context.Context) {
	service.refreshObservable = rxgo.Interval(rxgo.WithDuration(5 * time.Minute))
	service.refreshObservable.DoOnNext(func(i interface{}) {
		blacklist, err := service.repository.FindByCode(ctx, metadataKeyBlacklist)
		if err != nil {
			return
		}
		if service.cache == nil {
			service.cache = make(map[string]struct{})
		}
		for _, item := range blacklist {
			service.cache[item] = struct{}{}
		}
	})
}

func (service *SystemParameterServiceImpl) ShouldRecord(ctx context.Context, code string) bool {
	if service == nil {
		return true
	}
	_, ok := service.cache[code]
	return !ok
}

func (service *SystemParameterServiceImpl) Shutdown(ctx context.Context) error {
	if service.refreshObservable != nil {
		_, cancel := service.refreshObservable.Connect(ctx)
		cancel()
	}
	return service.repository.Shutdown(ctx)
}
