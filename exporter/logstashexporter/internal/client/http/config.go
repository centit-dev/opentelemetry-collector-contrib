package http // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/client/http

import (
	"fmt"
	"net/url"
)

type HttpConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c *HttpConfig) Validate() error {
	_, err := url.Parse(c.Endpoint)
	if err != nil {
		return fmt.Errorf("http endpoint %s invalid", c.Endpoint)
	}
	return nil
}
