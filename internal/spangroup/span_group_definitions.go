package spangroup

import (
	"strings"
)

const (
	opEqual              = "="
	opNotEqual           = "!="
	opGreaterThan        = ">"
	opGreaterThanOrEqual = ">="
	opLessThan           = "<"
	opLessThanOrEqual    = "<="
	opIn                 = "in"
	opNotIn              = "not-in"
	opContains           = "contains"
	opDoesNotContain     = "does-not-contain"
	opStartsWith         = "starts-with"
	opDoesNotStartWith   = "does-not-start-with"
	opExists             = "exists"
	opDoesNotExist       = "does-not-exist"
)

type SpanGroupDefinitions []SpanGroupDefinition

func (definitions *SpanGroupDefinitions) Match(attributes *map[string]interface{}) bool {
	for _, definition := range *definitions {
		if !definition.Match(attributes) {
			return false
		}
	}
	return true
}

type SpanGroupDefinition struct {
	Column string
	Op     string
	Value  GroupDefinitionValue
}

func CreateSpanGroupDefinition(column string, op string, value string) SpanGroupDefinition {
	return SpanGroupDefinition{
		Column: column,
		Op:     op,
		Value:  GroupDefinitionValue{StringValues: []string{value}},
	}
}

func (definition *SpanGroupDefinition) Match(attributes *map[string]interface{}) bool {
	value, exists := (*attributes)[definition.Column]
	switch definition.Op {
	// number or string
	case opEqual,
		opNotEqual,
		// number
		opGreaterThan,
		opGreaterThanOrEqual,
		opLessThan,
		opLessThanOrEqual,
		// string
		opContains,
		opDoesNotContain,
		opStartsWith,
		opDoesNotStartWith:
		return exists && definition.Value.compare(definition.Op, value)
	// array
	case opIn, opNotIn:
		if !exists {
			return false
		}
		switch raw := value.(type) {
		case float64:
			return arrayContain(definition.Op, definition.Value.numberValues, raw)
		case string:
			return arrayContain(definition.Op, definition.Value.StringValues, raw)
		}
		return false
	case opExists:
		return exists
	case opDoesNotExist:
		return !exists
	default:
		raw, ok := value.(bool)
		return ok && definition.Value.boolValue == raw
	}
}

// union types of float64, string, bool
type GroupDefinitionValue struct {
	numberValues []float64
	StringValues []string
	boolValue    bool
}

func CreateDefinitionValue(value interface{}) GroupDefinitionValue {
	if value == nil {
		return GroupDefinitionValue{}
	}
	switch raw := value.(type) {
	case float64:
		return GroupDefinitionValue{numberValues: []float64{raw}}
	case string:
		return GroupDefinitionValue{StringValues: []string{raw}}
	case bool:
		return GroupDefinitionValue{boolValue: value.(bool)}
	}

	// check if an array
	arr, isArray := value.([]interface{})
	if !isArray {
		return GroupDefinitionValue{}
	}

	// try to convert it to a number array or a string array
	numberArray := make([]float64, 0, len(arr))
	stringArray := make([]string, 0, len(arr))
	for _, item := range arr {
		// if the array contains both numbers and strings, it's invalid
		if len(numberArray) > 0 && len(stringArray) > 0 {
			return GroupDefinitionValue{}
		}
		switch rawItem := item.(type) {
		case float64:
			numberArray = append(numberArray, rawItem)
		case string:
			stringArray = append(stringArray, rawItem)
		default:
			// if the array contains other types, it's invalid
			return GroupDefinitionValue{}
		}
	}

	return GroupDefinitionValue{
		numberValues: numberArray,
		StringValues: stringArray,
	}
}

func (value *GroupDefinitionValue) compare(op string, actual interface{}) bool {
	switch raw := actual.(type) {
	case int64:
		return compareNumber(op, value.numberValues[0], float64(raw))
	case float64:
		return compareNumber(op, value.numberValues[0], raw)
	case string:
		return compareString(op, value.StringValues[0], raw)
	default:
		return false
	}
}

// actual comes from the attributes
// expected is the value in the definition
func compareNumber(op string, expected float64, actual float64) bool {
	switch op {
	case opEqual:
		return actual == expected
	case opNotEqual:
		return actual != expected
	case opGreaterThan:
		return actual > expected
	case opGreaterThanOrEqual:
		return actual >= expected
	case opLessThan:
		return actual < expected
	case opLessThanOrEqual:
		return actual <= expected
	default:
		return false
	}
}

// actual comes from the attributes
// expected is the value in the definition
func compareString(op string, expected string, actual string) bool {
	switch op {
	case opEqual:
		return actual == expected
	case opNotEqual:
		return actual != expected
	case opContains:
		return strings.Contains(actual, expected)
	case opDoesNotContain:
		return !strings.Contains(actual, expected)
	case opStartsWith:
		return strings.HasPrefix(actual, expected)
	case opDoesNotStartWith:
		return !strings.HasPrefix(actual, expected)
	default:
		return false
	}
}

// actual comes from the attributes
// expected is the value in the definition
func arrayContain(op string, expected interface{}, actual interface{}) bool {
	var contained bool

	if expectedNumbers, ok := expected.([]float64); ok {
		for _, item := range expectedNumbers {
			if item == actual {
				contained = true
				break
			}
		}
	} else if expectedStrings, ok := expected.([]string); ok {
		for _, item := range expectedStrings {
			if item == actual {
				contained = true
				break
			}
		}
	}
	if op == opIn {
		return contained
	}

	return !contained
}
