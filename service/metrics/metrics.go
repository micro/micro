// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/micro/v3/metrics/go

// Package metrics is for instrumentation and debugging
package metrics

import (
	"sync"
	"time"
)

// Tags is a map of fields to add to a metric:
type Tags map[string]string

// Reporter is an interface for collecting and instrumenting metrics
type Reporter interface {
	Count(id string, value int64, tags Tags) error
	Gauge(id string, value float64, tags Tags) error
	Timing(id string, value time.Duration, tags Tags) error
}

var (
	// DefaultMetricsReporter implementation
	DefaultMetricsReporter Reporter
	initialised            bool = false
	initialisedMutex       sync.Mutex
)

// IsSet lets you know if the DefaultMetricsReporter has been set already:
func IsSet() bool {
	initialisedMutex.Lock()
	defer initialisedMutex.Unlock()
	return initialised
}

// SetDefaultMetricsReporter allows other packages (such as profiles) to set the DefaultMetricsReporter
// The "initialised" flag prevents this from being overwritten (because other packages may already be using it)
func SetDefaultMetricsReporter(defaultReporter Reporter) {
	initialisedMutex.Lock()
	defer initialisedMutex.Unlock()
	if !initialised {
		DefaultMetricsReporter = defaultReporter
		initialised = true
	}
}

// Count submits a counter metric using the DefaultMetricsReporter:
func Count(id string, value int64, tags Tags) error {
	return DefaultMetricsReporter.Count(id, value, tags)
}

// Gauge submits a gauge metric using the DefaultMetricsReporter:
func Gauge(id string, value float64, tags Tags) error {
	return DefaultMetricsReporter.Gauge(id, value, tags)
}

// Timing submits a timing metric using the DefaultMetricsReporter:
func Timing(id string, value time.Duration, tags Tags) error {
	return DefaultMetricsReporter.Timing(id, value, tags)
}
