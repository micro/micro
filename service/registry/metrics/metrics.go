// Package metrics allows the profile package to set the reporter implementation
package metrics

import (
	"github.com/micro/go-micro/v3/metrics"
	"github.com/micro/go-micro/v3/metrics/noop"
)

var (
	// DefaultMetricsReporter implementation
	DefaultMetricsReporter metrics.Reporter = noop.New()
)
