// Package metrics allows the profile package to set the reporter implementation
package metrics

import (
	"time"

	"github.com/micro/go-micro/v3/metrics"
	"github.com/micro/go-micro/v3/metrics/noop"
	"github.com/micro/go-micro/v3/metrics/wrapper"
)

var (
	// DefaultMetricsReporter implementation
	DefaultMetricsReporter metrics.Reporter = noop.New()
)

// Count submits a counter metric using the DefaultMetricsReporter:
func Count(id string, value int64, tags metrics.Tags) error {
	return DefaultMetricsReporter.Count(id, value, tags)
}

// Gauge submits a gauge metric using the DefaultMetricsReporter:
func Gauge(id string, value float64, tags metrics.Tags) error {
	return DefaultMetricsReporter.Gauge(id, value, tags)
}

// Timing submits a timing metric using the DefaultMetricsReporter:
func Timing(id string, value time.Duration, tags metrics.Tags) error {
	return DefaultMetricsReporter.Timing(id, value, tags)
}

// Wrapper returns a metrics wrapper using the DefaultMetricsReporter:
func Wrapper() *wrapper.Wrapper {
	return wrapper.New(DefaultMetricsReporter)
}
