package httpbodyprocessor

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

const (
	httpRequestBodyKey     = "http.request.body"
	httpResponseBodyKey    = "http.response.body"
	httpReqContentTypeKey  = "http.request.header.content-type"
	httpRespContentTypeKey = "http.response.header.content-type"

	httpContentTypeXMLKey  = "text/xml"
	httpContentTypeJsonKey = "application/json"
)

const UnknowContentTypeBodyOutputLen = 5

type Processor struct {
	logger *zap.Logger
}

func NewProcessor(logger *zap.Logger) *Processor {
	return &Processor{logger}
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
	switch ParseHttpContentTypeByBody(value) {
	case httpContentTypeJsonKey:
		processJson(attributes, httpRequestBodyKey, &value)
	case httpContentTypeXMLKey:
		if err := processXml(attributes, httpRequestBodyKey, &value); err != nil {
			processor.logger.Error("process req xml fail", zap.Error(err))
			return
		}
	default:
		bodyOuput := value.Str()
		if bodyOuput != "" && len(bodyOuput) > UnknowContentTypeBodyOutputLen {
			bodyOuput = bodyOuput[:UnknowContentTypeBodyOutputLen] + "..."
		}
		processor.logger.Error("unreconize http req content type",
			zap.String("body", bodyOuput))
		return
	}
	attributes.Remove(httpRequestBodyKey)
}

func (processor *Processor) processResponseBody(attributes *pcommon.Map) {
	value, ok := attributes.Get(httpResponseBodyKey)
	if !ok {
		return
	}

	switch ParseHttpContentTypeByBody(value) {
	case httpContentTypeJsonKey:
		processJson(attributes, httpResponseBodyKey, &value)
	case httpContentTypeXMLKey:
		if err := processXml(attributes, httpResponseBodyKey, &value); err != nil {
			processor.logger.Error("process resp xml fail", zap.Error(err))
			return
		}
	default:
		bodyOuput := value.Str()
		if bodyOuput != "" && len(bodyOuput) > UnknowContentTypeBodyOutputLen {
			bodyOuput = bodyOuput[:UnknowContentTypeBodyOutputLen] + "..."
		}
		processor.logger.Error("unreconize http resp content type",
			zap.String("body", bodyOuput))
		return
	}
	attributes.Remove(httpResponseBodyKey)
}

func processXml(attributes *pcommon.Map, prefix string, value *pcommon.Value) error {
	xmlValue := value.Str()
	if xmlValue == "" {
		return nil
	}
	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlValue)))
	var root Node
	if err := decoder.Decode(&root); err != nil {
		return err
	}
	// 平坦化处理
	flatMap := make(map[string]string)
	flattenNode(prefix, root, flatMap)
	for k, v := range flatMap {
		attributes.PutStr(k, v)
	}
	return nil
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

type Node struct {
	XMLName xml.Name
	Content string     `xml:",chardata"`
	Nodes   []Node     `xml:",any"`
	Attr    []xml.Attr `xml:",attr"`
}

// flattenNode 递归地平坦化 XML 节点
func flattenNode(prefix string, node Node, flatMap map[string]string) {
	fullKey := node.XMLName.Local
	if prefix != "" {
		fullKey = prefix + "." + node.XMLName.Local
	}

	// 添加节点的属性
	for _, attr := range node.Attr {
		attrKey := fullKey + "." + attr.Name.Local
		flatMap[attrKey] = attr.Value
	}

	// 递归处理子节点
	for _, child := range node.Nodes {
		flattenNode(fullKey, child, flatMap)
	}

	// 如果当前节点是叶子节点，添加节点的文本内容
	if len(node.Nodes) == 0 && node.Content != "" {
		flatMap[fullKey] = node.Content
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

func ParseHttpContentTypeByBody(value pcommon.Value) string {
	content := value.Str()
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "{") {
		return httpContentTypeJsonKey
	} else if strings.HasPrefix(content, "<") {
		return httpContentTypeXMLKey
	} else {
		return ""
	}
}

func ParseHttpContentType(value pcommon.Value) string {
	switch value.Type() {
	case pcommon.ValueTypeStr:
		return value.Str()
	case pcommon.ValueTypeSlice:
		vs := value.Slice()
		if vs.Len() > 0 {
			return vs.At(0).Str()
		}
	}
	return ""
}
