// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package delaymonitorextension

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/delaymonitorextension/processor"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type delayMonitorExtension struct {
	config                    *Config
	telemetry                 component.TelemetrySettings
	delayMonitorSpanProcessor *processor.SpanProcessor
	handleDelayHistogram      metric.Float64Histogram
}

// registerableTracerProvider is a tracer that supports
// the SDK methods RegisterSpanProcessor and UnregisterSpanProcessor.
//
// We use an interface instead of casting to the SDK tracer type to support tracer providers
// that extend the SDK.
type registerableTracerProvider interface {
	// RegisterSpanProcessor adds the given SpanProcessor to the list of SpanProcessors.
	// https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProvider.RegisterSpanProcessor.
	RegisterSpanProcessor(SpanProcessor trace.SpanProcessor)

	// UnregisterSpanProcessor removes the given SpanProcessor from the list of SpanProcessors.
	// https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProvider.UnregisterSpanProcessor.
	UnregisterSpanProcessor(SpanProcessor trace.SpanProcessor)
}

func (dme *delayMonitorExtension) Start(_ context.Context, host component.Host) error {
	sdktracer, ok := dme.telemetry.TracerProvider.(registerableTracerProvider)
	if ok {
		sdktracer.RegisterSpanProcessor(dme.delayMonitorSpanProcessor)
		dme.telemetry.Logger.Info("Registered delay monitor span processor on tracer provider")
	} else {
		dme.telemetry.Logger.Warn("Delay monitor span processor registration is not available")
	}
	return nil
}

func (dme *delayMonitorExtension) Shutdown(context.Context) error {

	sdktracer, ok := dme.telemetry.TracerProvider.(registerableTracerProvider)
	if ok {
		sdktracer.UnregisterSpanProcessor(dme.delayMonitorSpanProcessor)
		dme.telemetry.Logger.Info("Unregistered delay monitor span processor on tracer provider")
	} else {
		dme.telemetry.Logger.Warn("Delay monitor span processor registration is not available")
	}

	return nil
}

const (
	extensionScope = "go.opentelemetry.io/collector/extension/extension"
)

func newServer(config *Config, telemetry component.TelemetrySettings) (*delayMonitorExtension, error) {
	delayM, err := telemetry.MeterProvider.Meter(extensionScope).Float64Histogram("component_handle_delay",
		metric.WithExplicitBucketBoundaries([]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5}...))

	if err != nil {
		return nil, err
	}
	return &delayMonitorExtension{
		config:                    config,
		telemetry:                 telemetry,
		delayMonitorSpanProcessor: processor.NewSpanProcessor(delayM),
		handleDelayHistogram:      delayM,
	}, nil
}
