package logstashexporter

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/client/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"path/filepath"
	"testing"
)

func TestHttpConfig_Validate(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata/http", "config.yaml"))
	require.NoError(t, err)

	type fields struct {
		Type    string
		httpCfg http.HttpConfig
	}
	tests := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id: component.NewIDWithName(Type, ""),
			expected: &Config{
				Type: "http",
				HttpCfg: http.HttpConfig{
					Endpoint: "http://127.0.0.1:1234",
					Username: "admin",
					Password: "123456",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			c := &Config{}
			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, c))

			assert.NoError(t, component.ValidateConfig(c))
			assert.Equal(t, tt.expected, c)
		})
	}
}
