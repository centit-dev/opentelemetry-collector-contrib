package http // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/client/http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/scheme"
)

const (
	jsonContentType = "application/json"

	statusCodeBackpressure = 429
	statusCodeInvalidAuth  = 401
	statusCodeSuccess      = 200
)

type LsHttpClient struct {
	cfg    *HttpConfig
	client *http.Client
	auth   string
}

func NewLsHttpClient(cfg *HttpConfig) *LsHttpClient {
	encodedAuth := ""
	cli, _ := InitHttpClient(cfg)
	if cfg.Username != "" && cfg.Password != "" {
		auth := cfg.Username + ":" + cfg.Password
		encodedAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	}
	return &LsHttpClient{cfg: cfg, client: cli, auth: encodedAuth}
}

func (c *LsHttpClient) getBaseHeader() http.Header {
	header := http.Header{}
	header.Add("Content-Type", jsonContentType)
	if c.auth != "" {
		header.Add("Authorization", c.auth)
	}
	return header
}

func (c *LsHttpClient) Handle(ctx context.Context, msg *scheme.LogsOutput) (err error) {
	formatMsg, err := json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("format output fail, err: %v", err)
		return
	}
	req, err := http.NewRequest("POST", c.cfg.Endpoint, bytes.NewBuffer(formatMsg))
	if err != nil {
		err = fmt.Errorf("make request fail, err: %v", err)
		return
	}
	req = req.WithContext(ctx)
	req.Header = c.getBaseHeader()
	resp, err := c.client.Do(req)
	if err != nil {
		err = fmt.Errorf("request logstash fail, err: %v", err)
		return
	}
	// todo check http resp code
	switch resp.StatusCode {
	case statusCodeSuccess:
		err = nil
	case statusCodeBackpressure:
		err = fmt.Errorf("logstash backend busy")
	case statusCodeInvalidAuth:
		err = fmt.Errorf("credential invalid")
	default:
		err = fmt.Errorf("internal error, status code: %d", resp.StatusCode)
	}
	return
}

func InitHttpClient(cfg *HttpConfig) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	clientTransport := (http.RoundTripper)(transport)
	return &http.Client{
		Transport: clientTransport,
	}, nil
}
