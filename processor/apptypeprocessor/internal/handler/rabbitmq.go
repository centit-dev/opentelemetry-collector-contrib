package handler

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

const (
	NetworkPeerAddressAttrName = "network.peer.address"
	NetworkPeerPortAttrName    = "network.peer.port"
)

type RabbitmqHandler struct {
}

func (h *RabbitmqHandler) GetHostPort(attrs *pcommon.Map) (host string, port string) {
	if ho, ok := attrs.Get(NetworkPeerAddressAttrName); ok {
		host = ho.Str()
	}
	if po, ok := attrs.Get(NetworkPeerPortAttrName); ok {
		port = fmt.Sprintf("%d", po.Int())
	}
	return
}
