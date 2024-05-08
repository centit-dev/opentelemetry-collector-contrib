package logstashexporter

import (
	"fmt"
	"net/url"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/client/http"
)

type Config struct {
	DryRun  bool            `mapstructure:"dryrun"`
	Type    string          `mapstructure:"type"`
	HttpCfg http.HttpConfig `mapstructure:"http_config"`
}

func (c *Config) Validate() error {
	if c.DryRun {
		return nil
	}
	if err := c.HttpCfg.Validate(); err != nil {
		return err
	}
	_, err := url.Parse(c.HttpCfg.Endpoint)
	if err != nil {
		return fmt.Errorf("http endpoint %s invalid", c.HttpCfg.Endpoint)
	}
	return nil
}
