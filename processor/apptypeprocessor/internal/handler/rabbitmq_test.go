package handler

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"testing"
)

func TestRabbitmqHandler_GetHostPort(t *testing.T) {
	attr := pcommon.NewMap()
	attr.PutStr(NetworkPeerAddressAttrName, "127.0.0.1")
	attr.PutInt(NetworkPeerPortAttrName, 5678)
	type args struct {
		attrs *pcommon.Map
	}
	tests := []struct {
		name     string
		args     args
		wantHost string
		wantPort string
	}{
		{
			name: "normal",
			args: args{
				attrs: &attr,
			},
			wantHost: "127.0.0.1",
			wantPort: "5678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &RabbitmqHandler{}
			gotHost, gotPort := h.GetHostPort(tt.args.attrs)
			if gotHost != tt.wantHost {
				t.Errorf("GetHostPort() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotPort != tt.wantPort {
				t.Errorf("GetHostPort() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}
