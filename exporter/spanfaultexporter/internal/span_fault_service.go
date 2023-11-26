package internal

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
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

type spanFaultOperator int

const (
	create spanFaultOperator = iota
	update
)

type spanFaultEntry struct {
	op   spanFaultOperator
	item *ent.SpanFault
}

type spanTree struct {
	rootSpan *ent.SpanFault
	spans    map[string]*spanTreeItem
	children map[string][]*spanTreeItem
}

type spanTreeItem struct {
	span              *ent.SpanFault
	hasRootCauseChild bool
}

type SpanFaultService interface {
	Start(ctx context.Context)
	Save(ctx context.Context, cause *ent.SpanFault) error
	Shutdown(ctx context.Context) error
}

type SpanFaultServiceImpl struct {
	logger *zap.Logger
	repo   SpanFaultRepository
	cache  *expirable.LRU[string, *spanTree]

	causeChannel           chan rxgo.Item
	observable             rxgo.Observable
	batchSize              int
	intervalInMilliseconds int
}

func CreateSpanFaultService(cacheConfig *CacheConfig, batchConfig *BatchConfig, repo SpanFaultRepository, logger *zap.Logger) SpanFaultService {
	channel := make(chan rxgo.Item, 1000)
	cache := expirable.NewLRU[string, *spanTree](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	return &SpanFaultServiceImpl{
		logger:                 logger,
		repo:                   repo,
		cache:                  cache,
		causeChannel:           channel,
		observable:             rxgo.FromChannel(channel),
		batchSize:              batchConfig.BatchSize,
		intervalInMilliseconds: batchConfig.IntervalInMilliseconds,
	}
}

func (service *SpanFaultServiceImpl) Start(ctx context.Context) {
	service.observable.
		BufferWithTimeOrCount(
			rxgo.WithDuration(time.Millisecond*time.Duration(service.intervalInMilliseconds)),
			service.batchSize).
		DoOnNext(func(items interface{}) {
			values := items.([]interface{})
			seenFaults := make(map[string]bool, len(values))
			creates := make([]*ent.SpanFault, 0, len(values))
			updates := make([]*ent.SpanFault, 0, len(values))
			for _, value := range values {
				entry, ok := value.(*spanFaultEntry)
				if !ok {
					continue
				}
				// remove duplicated creates/updates and only keep the first one
				// because all are just pointers to the same object
				// the order of the incoming items is guaranteed by the go channel
				// so it always prefer creating over updating
				if _, ok := seenFaults[entry.item.ID]; ok {
					continue
				}
				seenFaults[entry.item.ID] = true
				switch entry.op {
				case create:
					creates = append(creates, entry.item)
				default:
					updates = append(updates, entry.item)
				}
			}
			err := service.repo.SaveAll(ctx, creates, updates)
			if err != nil {
				service.logger.Error("error saving span faults", zap.Error(err))
			}
		}, rxgo.WithPool(10))
}

func (service *SpanFaultServiceImpl) Save(ctx context.Context, span *ent.SpanFault) error {
	// get or cache the tree
	tree, err := service.getOrCacheTree(ctx, span.TraceId)
	if err != nil {
		return err
	}

	// update the tree and the channel
	service.addSpan(tree, span)
	return nil
}

func (service *SpanFaultServiceImpl) getOrCacheTree(ctx context.Context, traceId string) (*spanTree, error) {
	tree, ok := service.cache.Get(traceId)
	if ok {
		return tree, nil
	}
	tree = &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	service.cache.Add(traceId, tree)
	spans, err := service.repo.GetSpanFaultsByTraceId(ctx, traceId)
	if err != nil {
		return nil, err
	}
	for _, span := range spans {
		service.addSpan(tree, span)
	}
	return tree, nil
}

// addSpan have 6 cases:
// 1. add the first span: add as parents and children
// 2. add the second span: add as parents and children
// 3. add as a child:
//   - if the child is a root cause, *update* the parent
//   - if the child is not a root cause, the parent is not changed
//
// 4. add as a parent:
//   - if the child is a root cause or has a root cause child, the parent is not a root cause and has a root cause child
//   - if the child is not a root cause and has no root cause child, the parent is a root cause and has no root cause child
//
// 5. add as a child having children:
//   - if the grand children have a root cause or has a root cause grand grand child, the child is not a root cause and has a root cause child
//   - if the grand children have no root cause and has no root cause grand grand child, the child is a root cause and has no root cause child
//   - if the child is a root cause or has a root cause child, *update* the parent if the parent is a root cause
//   - if the child is not a root cause and has no root cause child, the parent is a root cause and has no root cause child
//
// 6. add as root span: update all spans with root info
//
// only the direct parent of the new span is possibly updated
func (service *SpanFaultServiceImpl) addSpan(tree *spanTree, span *ent.SpanFault) {
	// find or create the root span
	if span.ParentSpanId == "" {
		tree.rootSpan = span
		span.RootServiceName = span.ServiceName
		span.RootSpanName = span.SpanName
		for _, item := range tree.spans {
			item.span.RootServiceName = span.ServiceName
			item.span.RootSpanName = span.SpanName
			service.causeChannel <- rxgo.Of(&spanFaultEntry{update, item.span})
		}
	} else if tree.rootSpan != nil {
		span.RootServiceName = tree.rootSpan.ServiceName
		span.RootSpanName = tree.rootSpan.SpanName
	}
	item := &spanTreeItem{span, false}

	// add as parents
	tree.spans[item.span.ID] = item
	// add as children
	if children, ok := tree.children[item.span.ParentSpanId]; ok {
		tree.children[item.span.ParentSpanId] = append(children, item)
	} else {
		tree.children[item.span.ParentSpanId] = []*spanTreeItem{item}
	}

	// update the root cause
	if item.span.FaultKind != "" {
		item.span.IsRoot = true
	}

	// update the current span if the children have a root cause
	children, hasChild := tree.children[item.span.ID]
	if hasChild {
		for _, child := range children {
			if child.hasRootCauseChild || child.span.IsRoot {
				item.span.IsRoot = false
				item.hasRootCauseChild = true
				break
			}
		}
	}

	// update the parent if the parent has a root cause child
	parent, hasParent := tree.spans[item.span.ParentSpanId]
	if hasParent && (item.hasRootCauseChild || item.span.IsRoot) {
		parent.hasRootCauseChild = true
		if parent.span.IsRoot {
			parent.span.IsRoot = false
			service.causeChannel <- rxgo.Of(&spanFaultEntry{update, parent.span})
		}
	}

	// always create the new span fault
	service.causeChannel <- rxgo.Of(&spanFaultEntry{create, span})
}

func (service *SpanFaultServiceImpl) Shutdown(ctx context.Context) error {
	_, cancel := service.observable.Connect(ctx)
	cancel()
	close(service.causeChannel)
	return service.repo.Shutdown(ctx)
}
