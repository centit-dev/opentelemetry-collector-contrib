// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package delaymonitorextension

import (
	"go.opentelemetry.io/collector/component"
)

type Config struct {
}

var _ component.Config = (*Config)(nil)

func (cfg *Config) Validate() error {
	return nil
}
