package spangroup

// an exception category can be found by multiple keys and the exception full name
type SpanGroup struct {
	value map[*SpanGroupDefinitions][]string
}

func CreateSpanGroup(data map[*SpanGroupDefinitions]string) *SpanGroup {
	m := make(map[*SpanGroupDefinitions][]string)
	for keys, value := range data {
		m[keys] = append(m[keys], value)
	}
	return &SpanGroup{m}
}

func (m *SpanGroup) Get(attributes *map[string]interface{}) []string {
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

func (m *SpanGroup) IsEmpty() bool {
	if m == nil {
		return true
	}
	return len(m.value) == 0
}
