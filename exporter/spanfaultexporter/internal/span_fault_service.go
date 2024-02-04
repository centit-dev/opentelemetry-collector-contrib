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
	Save(ctx context.Context, span []*ent.SpanFault) error
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

			err := service.repo.SaveAll(ctx, creates, []*ent.SpanFault{})
			if err != nil {
				service.logger.Error("failed to save batch", zap.Error(err))
			}
		}, rxgo.WithPool(10))
}

func (service *SpanFaultServiceImpl) Save(ctx context.Context, spans []*ent.SpanFault) error {
	// create the new trees
	trees := make(map[string]*spanTree)
	for _, span := range spans {
		tree, ok := trees[span.TraceId]
		if !ok {
			tree = &spanTree{
				spans:    make(map[string]*spanTreeItem),
				children: make(map[string][]*spanTreeItem),
			}
			trees[span.TraceId] = tree
			continue
		}

		service.addSpan(tree, span)
	}

	// send the new trees to the channel
	for _, tree := range trees {
		for _, item := range tree.spans {
			if !item.span.IsCause {
				continue
			}
			if tree.rootSpan == nil {
				break
			}
			tree.rootSpan.FaultKind = item.span.FaultKind
			// try to conclude the fault as the business fault
			// but if a system fault comes up, the final fault is system fault
			// TODO remove hard-coding fault kind
			if item.span.FaultKind == "SystemFault" {
				break
			}
		}
		if tree.rootSpan != nil {
			service.faultChannel <- rxgo.Of(tree.rootSpan)
		}
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
		item.span.IsCause = true
	}

	// update the current span if the children have a root cause
	children, hasChild := tree.children[item.span.ID]
	if hasChild {
		for _, child := range children {
			if child.hasFaultChild || child.span.IsCause {
				item.span.IsCause = false
				item.hasFaultChild = true
				break
			}
		}
	}

	// update the parent if the parent has a root cause child
	parent, hasParent := tree.spans[item.span.ParentSpanId]
	if hasParent && (item.hasFaultChild || item.span.IsCause) {
		parent.hasFaultChild = true
		if parent.span.IsCause {
			parent.span.IsCause = false
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
	return service.repo.Shutdown(ctx)
}
