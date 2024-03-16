package internal

import (
	"context"
	"testing"
	"time"

	"github.com/reactivex/rxgo/v2"
)

func TestStart(t *testing.T) {
	client := createTestClient()
	repository := CreateSystemParameterRepository(client)
	service := CreateSystemParameterService(repository)
	service.refreshObservable = rxgo.Interval(rxgo.WithDuration(1 * time.Second))
	ctx := context.Background()
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

	time.Sleep(6 * time.Minute)
}
