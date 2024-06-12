package internal

import (
	"context"
	"sort"
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
	childrenGroup := make(map[string][]*spanTreeItem)

	for _, item := range items {
		if item.ParentSpanId != "" {
			if siblings, exists := childrenGroup[item.ParentSpanId]; exists {
				childrenGroup[item.ParentSpanId] = append(siblings, item)
			} else {
				childrenGroup[item.ParentSpanId] = []*spanTreeItem{item}
			}
		}

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
		var cause *spanTreeItem
		var depth int16 = -1
		for _, span := range tree.spans {
			spanDepth := span.Depth(&tree.spans)
			if spanDepth > depth && span.FaultKind != "" {
				depth = spanDepth
				cause = span
			}
		}
		if tree.rootSpan == nil {
			continue
		}
		if cause == nil {
			cause = tree.rootSpan
		}
		cause.PlatformName = tree.rootSpan.PlatformName
		cause.AppCluster = tree.rootSpan.AppCluster
		cause.InstanceName = tree.rootSpan.InstanceName
		cause.RootServiceName = tree.rootSpan.ServiceName
		cause.RootSpanName = tree.rootSpan.SpanName
		cause.RootDuration = tree.rootSpan.duration
		if cause.ParentSpanId != "" {
			parent, ok := tree.spans[cause.ParentSpanId]
			if ok {
				cause.Gap = cause.Timestamp.Sub(parent.Timestamp).Nanoseconds()
			}
		}
		if children, exists := childrenGroup[cause.SpanId]; exists {
			cause.SelfDuration = service.calculateSelfDuration(cause, children)
		} else {
			cause.SelfDuration = cause.duration
		}
		service.faultChannel <- rxgo.Of(cause.SpanFault)
	}
	return nil
}

func (service *SpanFaultServiceImpl) calculateSelfDuration(span *spanTreeItem, children []*spanTreeItem) int64 {
	if len(children) == 0 {
		return span.duration
	}

	sort.Slice(children, func(i, j int) bool {
		return children[i].Timestamp.Before(children[j].Timestamp)
	})
	maxEndTime := span.Timestamp.Add(time.Duration(span.duration))
	intervals := make([][]int64, 0, len(children))
	for _, child := range children {
		if len(intervals) == 0 {
			intervals = append(intervals, []int64{child.Timestamp.UnixNano(), child.Timestamp.UnixNano() + child.duration, child.duration})
			continue
		}

		last_entry := intervals[len(intervals)-1]
		last_start := last_entry[0]
		last_end := last_entry[1]

		current_end := child.Timestamp.UnixNano() + child.duration
		if child.Timestamp.UnixNano() > last_end {
			// not overlap:
			// last:  |---------|
			// next:                |---------|
			intervals = append(intervals, []int64{child.Timestamp.UnixNano(), current_end, child.duration})
		} else if current_end > last_end {
			// overlaps:
			// last:  |---------|
			// next:       |---------|
			intervals[len(intervals)-1] = []int64{last_start, current_end, current_end - last_start}
		}
		// sorting by start time ascending ganrantees start_time > last_start_time
		// so this is impossible
		// last:         |---------|
		// next: |---------|

		// lastly, if the last end time is greater than the parent end time, we need to adjust the last interval
		// last:            |---------|
		// parent: |---------|
		last_entry = intervals[len(intervals)-1]
		last_start = last_entry[0]
		last_end = last_entry[1]
		if last_end > maxEndTime.UnixNano() {
			intervals[len(intervals)-1] = []int64{last_start, maxEndTime.UnixNano(), maxEndTime.UnixNano() - last_start}
			// and there is no point to continue
			break
		}
	}
	selfDuration := span.duration
	for _, interval := range intervals {
		selfDuration -= interval[2]
	}
	// client error may cause negative self duration
	if selfDuration < 0 {
		selfDuration = 0
	}
	return selfDuration
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
