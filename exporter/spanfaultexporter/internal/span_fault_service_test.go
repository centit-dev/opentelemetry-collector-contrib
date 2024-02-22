package internal

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSpanTreeItemDepthWithRootSpan(t *testing.T) {
	spans := make(map[string]*spanTreeItem)
	item := buildSpanFault(uuid.NewString())
	item.ParentSpanId = ""
	depth := item.Depth(&spans)
	assert.Equal(t, int16(0), depth)
}

func TestSpanTreeItemDepthWithIncompleteSpan(t *testing.T) {
	spans := make(map[string]*spanTreeItem)
	item := buildSpanFault(uuid.NewString())
	depth := item.Depth(&spans)
	assert.Equal(t, int16(-15_000), depth)
}

func TestSpanTreeItemDepthWithParent(t *testing.T) {
	parent := buildSpanFault(uuid.NewString())
	parent.ParentSpanId = ""
	item := buildSpanFault(uuid.NewString())
	item.ParentSpanId = "parent"
	spans := map[string]*spanTreeItem{
		"parent": parent,
	}
	depth := item.Depth(&spans)
	assert.Equal(t, int16(1), depth)
}
