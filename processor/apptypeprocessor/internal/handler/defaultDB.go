package handler

import (
	"net/url"

	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.22.0"
)

type DefaultDBHandler struct {
}

func (h *DefaultDBHandler) GetHostPort(attrs *pcommon.Map) (host string, port string) {
	if raw, ok := attrs.Get(conventions.AttributeDBConnectionString); ok {
		u, err := url.Parse(raw.Str())
		if err != nil {
			return
		}

		host = u.Hostname()
		port = u.Port()
	}
	return
}
