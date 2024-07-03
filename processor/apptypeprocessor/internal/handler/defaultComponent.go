package handler

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

type DefaultComponentHandler struct {
	UrlKey string
}

func (h *DefaultComponentHandler) GetHostPort(attrs *pcommon.Map) (host string, port string) {
	if url, ok := attrs.Get(h.UrlKey); ok {
		if url.Str() == "" {
			return
		}
		info := strings.Split(url.Str(), ":")
		if len(info) > 0 {
			host = info[0]
		}
		if len(info) > 1 {
			port = info[1]
		}
	}
	return
}
