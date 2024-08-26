package handler

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"net"
	"strings"
)

type ActivemqHandler struct {
}

func (h *ActivemqHandler) GetHostPort(attrs *pcommon.Map) (host string, port string) {
	if url, ok := attrs.Get(ActivemqConnectionAttrName); ok {
		trimmedInput := strings.TrimPrefix(url.Str(), "tcp://")
		host, port, _ = net.SplitHostPort(trimmedInput)
	}
	return
}
