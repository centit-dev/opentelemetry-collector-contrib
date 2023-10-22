package client

// an exception category can be found by multiple keys and the exception full name
type ExceptionCategoryMap struct {
	value map[*ExceptionCategoryDefinitions][]string
}

func CreateExceptionCategoryMap(data map[*ExceptionCategoryDefinitions]string) *ExceptionCategoryMap {
	m := make(map[*ExceptionCategoryDefinitions][]string)
	for keys, value := range data {
		m[keys] = append(m[keys], value)
	}
	return &ExceptionCategoryMap{m}
}

func (m *ExceptionCategoryMap) Get(attributes *map[string]interface{}) []string {
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

func (m *ExceptionCategoryMap) IsEmpty() bool {
	if m == nil {
		return true
	}
	return len(m.value) == 0
}
