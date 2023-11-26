package internal

import (
	"testing"

	"github.com/reactivex/rxgo/v2"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
)

func TestSpanFaultServiceImpl_addSpan_firstSpan(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span := buildSpanFault("trace1")
	service.addSpan(tree, span)

	_, ok := tree.children[span.ParentSpanId]
	if !ok {
		t.Errorf("failed to add span as children")
	}

	_, ok = tree.spans[span.ID]
	if !ok {
		t.Errorf("failed to add span as parents")
	}
}

func TestSpanFaultServiceImpl_addSpan_secondSpan(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	service.addSpan(tree, span2)

	_, ok := tree.children[span2.ParentSpanId]
	if !ok {
		t.Errorf("failed to add span as children")
	}

	_, ok = tree.spans[span2.ID]
	if !ok {
		t.Errorf("failed to add span as parents")
	}
}

func TestSpanFaultServiceImpl_addSpan_aChildSpan(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	span2.ParentSpanId = span1.ID
	service.addSpan(tree, span2)

	parent := tree.spans[span1.ID]
	if !parent.hasRootCauseChild {
		t.Errorf("the parent should have a root cause child")
	}
}

func TestSpanFaultServiceImpl_addSpan_aParentSpan(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	span2.ID = span1.ParentSpanId
	service.addSpan(tree, span2)

	parent := tree.spans[span2.ID]
	if !parent.hasRootCauseChild {
		t.Errorf("the parent should have a root cause child")
	}
}

func TestSpanFaultServiceImpl_addSpan_aChildSpanWithoutRootCause(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	span2.ParentSpanId = span1.ID
	span2.FaultKind = ""
	span2.IsRoot = false
	service.addSpan(tree, span2)

	parent := tree.spans[span1.ID]
	if parent.hasRootCauseChild || !parent.span.IsRoot {
		t.Errorf("the parent is a root cause")
	}
}

func TestSpanFaultServiceImpl_addSpan_noRootCauseUntilAddingAParent(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	span1.FaultKind = ""
	span1.IsRoot = false
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	span2.ID = span1.ParentSpanId
	service.addSpan(tree, span2)

	parent := tree.spans[span2.ID]
	if !parent.span.IsRoot {
		t.Errorf("the parent should be a root cause")
	}
}

func TestSpanFaultServiceImpl_addSpan_noRootCauseUntilAddingAChild(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	span1.FaultKind = ""
	span1.IsRoot = false
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	span2.ParentSpanId = span1.ID
	service.addSpan(tree, span2)

	parent := tree.spans[span1.ID]
	if !parent.hasRootCauseChild {
		t.Errorf("the parent should have a root cause child")
	}
}

func TestSpanFaultServiceImpl_addSpan_addANonFaultChildBetweenTwoFaults(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	service.addSpan(tree, span2)

	if !span1.IsRoot {
		t.Errorf("the first span should be a root cause")
	}

	if !span2.IsRoot {
		t.Errorf("the second span should be a root cause")
	}

	span3 := buildSpanFault("trace1")
	span3.ParentSpanId = span1.ID
	span3.ID = span2.ParentSpanId
	span3.FaultKind = ""
	span3.IsRoot = false
	service.addSpan(tree, span3)

	parent := tree.spans[span3.ParentSpanId]
	if !parent.hasRootCauseChild || parent.span.IsRoot {
		t.Errorf("the parent should have a root cause child")
	}

	child := tree.spans[span3.ID]
	if !child.hasRootCauseChild {
		t.Errorf("the child should have a root cause child")
	}

	grandchild := tree.spans[span2.ID]
	if !grandchild.span.IsRoot {
		t.Errorf("the grandchild should be a root cause")
	}
}

func TestSpanFaultServiceImpl_addSpan_addRootSpanAfterReceiving2(t *testing.T) {
	service := &SpanFaultServiceImpl{
		causeChannel: make(chan rxgo.Item, 100),
	}

	tree := &spanTree{
		spans:    make(map[string]*spanTreeItem),
		children: make(map[string][]*spanTreeItem),
	}
	span1 := buildSpanFault("trace1")
	service.addSpan(tree, span1)

	span2 := buildSpanFault("trace1")
	service.addSpan(tree, span2)

	span3 := buildSpanFault("trace1")
	span3.ParentSpanId = ""
	service.addSpan(tree, span3)

	span4 := buildSpanFault("trace1")
	service.addSpan(tree, span4)

	if span1.RootServiceName != span3.ServiceName {
		t.Errorf("the span should be updated with the root service name")
	}

	ops := make([]spanFaultOperator, 0, 6)
	spans := make([]*ent.SpanFault, 0, 6)
	for item := range service.causeChannel {
		entry := item.V.(*spanFaultEntry)
		ops = append(ops, entry.op)
		spans = append(spans, entry.item)
		if len(ops) == 6 {
			close(service.causeChannel)
		}
	}

	if !compareSlices(ops, []spanFaultOperator{create, create, update, update, create, create}) {
		t.Errorf("the ops doesn't match")
	}

	for _, span := range spans {
		if span.RootServiceName == "" {
			t.Errorf("the span should be updated with the root service name")
		}
	}
}

func compareSlices(a, b []spanFaultOperator) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
