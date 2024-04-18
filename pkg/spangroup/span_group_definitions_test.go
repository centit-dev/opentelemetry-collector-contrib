package spangroup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpanGroupDefinitionsMatch(t *testing.T) {
	definitions := SpanGroupDefinitions{
		{
			Column: "span.kind",
			Op:     "!=",
			Value:  GroupDefinitionValue{StringValues: []string{""}},
		},
		{
			Column: "span.kind",
			Op:     "!=",
			Value:  GroupDefinitionValue{StringValues: []string{"value"}},
		},
	}
	value := &map[string]interface{}{"span.kind": "value2"}
	assert.True(t, definitions.Match(value))
}

func TestSpanGroupDefinitionMatch(t *testing.T) {
	gt := &SpanGroupDefinition{
		Column: "span.kind",
		Op:     ">",
		Value:  GroupDefinitionValue{numberValues: []float64{1.}},
	}
	value := &map[string]interface{}{"span.kind": 2.}
	assert.True(t, gt.Match(value))

	eq := &SpanGroupDefinition{
		Column: "span.kind",
		Op:     "=",
		Value:  GroupDefinitionValue{StringValues: []string{"value"}},
	}
	value = &map[string]interface{}{"span.kind": "value"}
	assert.True(t, eq.Match(value))
}

func TestCreateDefinitionValue(t *testing.T) {
	stringValue := CreateDefinitionValue("value")
	assert.Equal(t, stringValue, GroupDefinitionValue{StringValues: []string{"value"}})

	numberValue := CreateDefinitionValue(1.)
	assert.Equal(t, numberValue, GroupDefinitionValue{numberValues: []float64{1.}})

	boolValue := CreateDefinitionValue(true)
	assert.Equal(t, boolValue, GroupDefinitionValue{boolValue: true})
}

func TestArrayContain(t *testing.T) {
	assert.True(t, arrayContain("in", []float64{1., 2.}, 2.))
	assert.True(t, arrayContain("not-in", []float64{1., 2.}, 3.))
}
