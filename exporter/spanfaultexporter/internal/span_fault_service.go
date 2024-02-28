package internal

import (
	"context"
	"time"

	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
	"go.uber.org/zap"
)

type CacheConfig struct {
	MaxSize         int `mapstructure:"max_size"`
	ExpireInMinutes int `mapstructure:"expire_in_minutes"`
}

type BatchConfig struct {
	BatchSize int `mapstructure:"batch_size"`
	// use millisecond for testing purpose and be more flexible
	IntervalInMilliseconds int `mapstructure:"interval_in_milliseconds"`
}

type spanTree struct {
	rootSpan *spanTreeItem
	spans    map[string]*spanTreeItem
}

type spanTreeItem struct {
	*ent.SpanFault
	duration int64
	parent   *spanTreeItem
	depth    int16
}

func (item *spanTreeItem) Depth(spans *map[string]*spanTreeItem) int16 {
	if item.depth > 0 {
		return item.depth
	}
	if item.ParentSpanId == "" {
		return 0
	}
	parent, ok := (*spans)[item.ParentSpanId]
	if !ok {
		if item.ParentSpanId == "" {
			return 0
		} else {
			return -15_000
		}
	}
	item.depth = parent.Depth(spans) + 1
	return item.depth
}

type SpanFaultService interface {
	Start(ctx context.Context)
	Save(ctx context.Context, span []*spanTreeItem) error
	Shutdown(ctx context.Context) error
}

type SpanFaultServiceImpl struct {
	logger *zap.Logger
	repo   SpanFaultRepository

	faultChannel           chan rxgo.Item
	faultObservable        rxgo.Observable
	batchSize              int
	intervalInMilliseconds int
}

func CreateSpanFaultService(batchConfig *BatchConfig, repo SpanFaultRepository, logger *zap.Logger) SpanFaultService {
	faultChannel := make(chan rxgo.Item, batchConfig.BatchSize+1000)
	return &SpanFaultServiceImpl{
		logger:                 logger,
		repo:                   repo,
		faultChannel:           faultChannel,
		faultObservable:        rxgo.FromChannel(faultChannel),
		batchSize:              batchConfig.BatchSize,
		intervalInMilliseconds: batchConfig.IntervalInMilliseconds,
	}
}

func (service *SpanFaultServiceImpl) Start(ctx context.Context) {
	service.faultObservable.
		BufferWithTimeOrCount(
			rxgo.WithDuration(time.Millisecond*time.Duration(service.intervalInMilliseconds)),
			service.batchSize).
		DoOnNext(func(items interface{}) {
			service.logger.Debug("batching")

			values := items.([]interface{})
			creates := make([]*ent.SpanFault, 0, len(values))
			for _, item := range values {
				span, ok := item.(*ent.SpanFault)
				if !ok {
					continue
				}

				creates = append(creates, span)
			}

			err := service.repo.SaveAll(ctx, creates)
			if err != nil {
				service.logger.Error("failed to save batch", zap.Error(err))
			}
		}, rxgo.WithPool(10))
}

func (service *SpanFaultServiceImpl) Save(ctx context.Context, items []*spanTreeItem) error {
	// create the new trees
	trees := make(map[string]*spanTree)
	for _, item := range items {
		tree, ok := trees[item.ID]
		if !ok {
			tree = &spanTree{
				spans: make(map[string]*spanTreeItem),
			}
			trees[item.ID] = tree
		}

		service.addSpan(tree, item)
	}

	for _, tree := range trees {
		var cause *ent.SpanFault
		var depth int16 = -1
		for _, span := range tree.spans {
			spanDepth := span.Depth(&tree.spans)
			if spanDepth > depth && span.FaultKind != "" {
				depth = spanDepth
				cause = span.SpanFault
			}
		}
		if tree.rootSpan == nil {
			continue
		}
		if cause == nil {
			cause = tree.rootSpan.SpanFault
		}
		cause.RootServiceName = tree.rootSpan.ServiceName
		cause.RootSpanName = tree.rootSpan.SpanName
		cause.RootDuration = tree.rootSpan.duration
		service.faultChannel <- rxgo.Of(cause)
	}
	return nil
}

// addSpan adds a span to the tree and find out the root span
func (service *SpanFaultServiceImpl) addSpan(tree *spanTree, item *spanTreeItem) {
	// find and set the root span
	if item.ParentSpanId == "" {
		tree.rootSpan = item
	}

	// add as parents
	tree.spans[item.SpanId] = item
}

func (service *SpanFaultServiceImpl) Shutdown(ctx context.Context) error {
	_, cancel := service.faultObservable.Connect(ctx)
	cancel()
	close(service.faultChannel)
	return service.repo.Shutdown(ctx)
}
