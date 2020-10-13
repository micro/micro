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
// Original source: github.com/micro/micro/v3/metrics/options.go

package metrics

var (
	// The Prometheus metrics will be made available on this port:
	defaultPrometheusListenAddress = ":9000"
	// This is the endpoint where the Prometheus metrics will be made available ("/metrics" is the default with Prometheus):
	defaultPath = "/metrics"
	// defaultPercentiles is the default spread of percentiles/quantiles we maintain for timings / histogram metrics:
	defaultPercentiles = []float64{0, 0.5, 0.75, 0.90, 0.95, 0.98, 0.99, 1}
)

// Option powers the configuration for metrics implementations:
type Option func(*Options)

// Options for metrics implementations:
type Options struct {
	Address     string
	DefaultTags Tags
	Path        string
	Percentiles []float64
}

// NewOptions prepares a set of options:
func NewOptions(opt ...Option) Options {
	opts := Options{
		Address:     defaultPrometheusListenAddress,
		DefaultTags: make(Tags),
		Path:        defaultPath,
		Percentiles: defaultPercentiles,
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// Path used to serve metrics over HTTP:
func Path(value string) Option {
	return func(o *Options) {
		o.Path = value
	}
}

// Address is the listen address to serve metrics on:
func Address(value string) Option {
	return func(o *Options) {
		o.Address = value
	}
}

// DefaultTags will be added to every metric:
func DefaultTags(value Tags) Option {
	return func(o *Options) {
		o.DefaultTags = value
	}
}

// Percentiles defines the desired spread of statistics for histogram / timing metrics:
func Percentiles(value []float64) Option {
	return func(o *Options) {
		o.Percentiles = value
	}
}
