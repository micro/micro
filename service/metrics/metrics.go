// Package metrics allows the profile package to set the reporter implementation
package metrics

import (
	"sync"
	"time"

	"github.com/micro/go-micro/v3/metrics"
	"github.com/micro/go-micro/v3/metrics/noop"
)

var (
	// DefaultMetricsReporter implementation
	DefaultMetricsReporter metrics.Reporter = noop.New()
	initialised            bool             = false
	initialisedMutex       sync.Mutex
)

// SetDefaultMetricsReporter allows other packages (such as profiles) to set the DefaultMetricsReporter
// The "initialised" flag prevents this from being overwritten (because other packages may already be using it)
func SetDefaultMetricsReporter(defaultReporter metrics.Reporter) {
	initialisedMutex.Lock()
	defer initialisedMutex.Unlock()
	if !initialised {
		DefaultMetricsReporter = defaultReporter
		initialised = true
	}
}

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
