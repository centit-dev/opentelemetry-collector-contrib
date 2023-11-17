package spangroup

// a group can have multiple groups of definitions
// each group of definitions have multiple definitions
// each definition has a column, an operator and a value that matches a list of spans
// this makes a group of definitions groups spans by requiring them to match all definitions in this group
// the grouped spans are labeled by the values of the map, eg, the group names
type SpanGroups struct {
	value map[*SpanGroupDefinitions][]string
}

func CreateSpanGroup(data map[*SpanGroupDefinitions]string) *SpanGroups {
	m := make(map[*SpanGroupDefinitions][]string)
	for keys, value := range data {
		m[keys] = append(m[keys], value)
	}
	return &SpanGroups{m}
}

func (m *SpanGroups) Get(attributes *map[string]interface{}) []string {
	matches := make(map[string]int)
	for definitions, values := range m.value {
		if definitions.Match(attributes) {
			for _, value := range values {
				matches[value] = 1
			}
		}
	}
	result := make([]string, 0, len(matches))
	for key := range matches {
		result = append(result, key)
	}
	return result
}

func (m *SpanGroups) IsEmpty() bool {
	if m == nil {
		return true
	}
	return len(m.value) == 0
}
