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
	span          *ent.SpanFault
	hasFaultChild bool
}

type SpanFaultService interface {
	Start(ctx context.Context)
	Save(ctx context.Context, spans []*ent.SpanFault) error
	Shutdown(ctx context.Context) error
}

type SpanFaultServiceImpl struct {
	logger      *zap.Logger
	repo        SpanFaultRepository
	cache       *expirable.LRU[string, *spanTree]
	ignoreCache *expirable.LRU[string, bool]

	faultChannel           chan rxgo.Item
	faultObservable        rxgo.Observable
	queueChannel           chan rxgo.Item
	queueObservable        rxgo.Observable
	batchSize              int
	intervalInMilliseconds int
}

func CreateSpanFaultService(cacheConfig *CacheConfig, batchConfig *BatchConfig, repo SpanFaultRepository, logger *zap.Logger) SpanFaultService {
	faultChannel := make(chan rxgo.Item, batchConfig.BatchSize+1000)
	queueChannel := make(chan rxgo.Item, cacheConfig.MaxSize+1000)
	cache := expirable.NewLRU[string, *spanTree](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	ignoreCache := expirable.NewLRU[string, bool](cacheConfig.MaxSize, nil, time.Minute*time.Duration(cacheConfig.ExpireInMinutes))
	return &SpanFaultServiceImpl{
		logger:                 logger,
		repo:                   repo,
		cache:                  cache,
		ignoreCache:            ignoreCache,
		faultChannel:           faultChannel,
		faultObservable:        rxgo.FromChannel(faultChannel),
		queueChannel:           queueChannel,
		queueObservable:        rxgo.FromChannel(queueChannel),
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

	// buffer by 1 second and wait for 9 seconds before processing
	// so the spans are processed after the trace is considered completed
	service.queueObservable.
		BufferWithTime(rxgo.WithDuration(time.Second)).
		DoOnNext(func(items interface{}) {
			timer := time.NewTimer(9 * time.Second)
			<-timer.C

			values := items.([]interface{})
			for _, value := range values {
				tree, ok := value.(*spanTree)
				if !ok {
					continue
				}
				for _, item := range tree.spans {
					if !item.span.IsRoot {
						continue
					}
					tree.rootSpan.FaultKind = item.span.FaultKind
					// TODO remove hard-coding fault kind
					if item.span.FaultKind == "SystemFault" {
						break
					}
				}
				service.cache.Remove(tree.rootSpan.TraceId)
				service.faultChannel <- rxgo.Of(&spanFaultEntry{create, tree.rootSpan})
			}
		}, rxgo.WithPool(15))
}

func (service *SpanFaultServiceImpl) Save(ctx context.Context, spans []*ent.SpanFault) error {
	// update the tree if exists
	traceIds := make(map[string]bool)
	for _, span := range spans {
		_, ok := service.cache.Get(span.TraceId)
		if ok {
			continue
		}
		traceIds[span.TraceId] = true
	}

	// or check if the trace is already created
	uniqueTraceIds := make([]string, 0, len(traceIds))
	for traceId := range traceIds {
		uniqueTraceIds = append(uniqueTraceIds, traceId)
	}
	faults, err := service.repo.GetSpanFaultsByTraceIds(ctx, uniqueTraceIds)
	if err != nil {
		return err
	}
	// find the new traces
	for _, fault := range faults {
		delete(traceIds, fault.TraceId)
	}

	// create the new trees
	for traceId := range traceIds {
		tree := &spanTree{
			spans:    make(map[string]*spanTreeItem),
			children: make(map[string][]*spanTreeItem),
		}
		service.cache.Add(traceId, tree)
		service.queueChannel <- rxgo.Of(tree)
	}

	// update the tree
	for _, span := range spans {
		tree, ok := service.cache.Get(span.TraceId)
		if !ok {
			continue
		}
		service.addSpan(tree, span)
	}
	return nil
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
			// service.causeChannel <- rxgo.Of(&spanFaultEntry{update, item.span})
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
			if child.hasFaultChild || child.span.IsRoot {
				item.span.IsRoot = false
				item.hasFaultChild = true
				break
			}
		}
	}

	// update the parent if the parent has a root cause child
	parent, hasParent := tree.spans[item.span.ParentSpanId]
	if hasParent && (item.hasFaultChild || item.span.IsRoot) {
		parent.hasFaultChild = true
		if parent.span.IsRoot {
			parent.span.IsRoot = false
			// service.causeChannel <- rxgo.Of(&spanFaultEntry{update, parent.span})
		}
	}

	// always create the new span fault
	// service.causeChannel <- rxgo.Of(&spanFaultEntry{create, span})
}

func (service *SpanFaultServiceImpl) Shutdown(ctx context.Context) error {
	_, cancel := service.faultObservable.Connect(ctx)
	cancel()
	close(service.faultChannel)

	_, cancel = service.queueObservable.Connect(ctx)
	cancel()
	close(service.queueChannel)
	return service.repo.Shutdown(ctx)
}
