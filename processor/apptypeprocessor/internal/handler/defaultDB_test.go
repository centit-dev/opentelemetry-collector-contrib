package handler

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"testing"
)

func TestDefaultDBHandler_GetHostPort(t *testing.T) {
	attr := pcommon.NewMap()
	attr.PutStr("db.connection_string", "127.0.0.1:11211")
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
			wantPort: "11211",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DefaultDBHandler{}
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
