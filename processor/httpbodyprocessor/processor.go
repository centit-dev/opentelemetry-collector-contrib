package httpbodyprocessor

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	httpRequestBodyKey  = "http.request.body"
	httpResponseBodyKey = "http.response.body"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (processor *Processor) ProcessTraces(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		resourceSpans := traces.ResourceSpans().At(i)
		for j := 0; j < resourceSpans.ScopeSpans().Len(); j++ {
			scopeSpans := resourceSpans.ScopeSpans().At(j)
			for k := 0; k < scopeSpans.Spans().Len(); k++ {
				span := scopeSpans.Spans().At(k)
				attributes := span.Attributes()
				processor.processHttpRequestBody(&attributes)
				processor.processResponseBody(&attributes)
			}
		}
	}
	return traces, nil
}

func (processor *Processor) processHttpRequestBody(attributes *pcommon.Map) {
	value, ok := attributes.Get(httpRequestBodyKey)
	if !ok {
		return
	}
	processJson(attributes, httpRequestBodyKey, &value)
	attributes.Remove(httpRequestBodyKey)
}

func (processor *Processor) processResponseBody(attributes *pcommon.Map) {
	value, ok := attributes.Get(httpResponseBodyKey)
	if !ok {
		return
	}
	processJson(attributes, httpResponseBodyKey, &value)
	attributes.Remove(httpResponseBodyKey)
}

func processJson(attributes *pcommon.Map, prefix string, value *pcommon.Value) {
	jsonValue := value.Str()
	if jsonValue == "" {
		return
	}
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(jsonValue), &input); err != nil {
		return
	}
	output := make(map[string]string)
	flattenJson(input, prefix, &output)
	for key, value := range output {
		attributes.PutStr(key, value)
	}
}

func flattenJson(input map[string]interface{}, parentKey string, output *map[string]string) {
	for key, value := range input {
		var currentKey string
		if parentKey == "" {
			currentKey = key
		} else {
			currentKey = fmt.Sprintf("%s.%s", parentKey, key)
		}
		switch value := value.(type) {
		case map[string]interface{}:
			flattenJson(value, currentKey, output)
		case string:
			(*output)[currentKey] = value
		default:
			if value == nil {
				continue
			}
			(*output)[currentKey] = fmt.Sprintf("%v", value)
		}
	}
}
