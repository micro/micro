// Copyright 2020 Asim Aslam
//
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
// Original source: github.com/micro/go-plugins/v3/metrics/prometheus/reporter.go

package prometheus

import (
	"net/http"
	"strings"

	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// quantileThresholds maps quantiles / percentiles to error thresholds (required by the Prometheus client).
	// Must be from our pre-defined set [0.0, 0.5, 0.75, 0.90, 0.95, 0.98, 0.99, 1]:
	quantileThresholds = map[float64]float64{0.0: 0, 0.5: 0.05, 0.75: 0.04, 0.90: 0.03, 0.95: 0.02, 0.98: 0.001, 1: 0}
)

// Reporter is an implementation of metrics.Reporter:
type Reporter struct {
	options            metrics.Options
	prometheusRegistry *prometheus.Registry
	metrics            metricFamily
}

// New returns a configured prometheus reporter:
func New(opts ...metrics.Option) (*Reporter, error) {
	options := metrics.NewOptions(opts...)

	// Make a prometheus registry (this keeps track of any metrics we generate):
	prometheusRegistry := prometheus.NewRegistry()
	prometheusRegistry.Register(prometheus.NewGoCollector())
	prometheusRegistry.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{Namespace: "goruntime"}))

	// Make a new Reporter:
	newReporter := &Reporter{
		options:            options,
		prometheusRegistry: prometheusRegistry,
	}

	// Add metrics families for each type:
	newReporter.metrics = newReporter.newMetricFamily()

	// Handle the metrics endpoint with prometheus:
	log.Infof("Metrics/Prometheus [http] Listening on %s%s", options.Address, options.Path)
	http.Handle(options.Path, promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	go http.ListenAndServe(options.Address, nil)

	return newReporter, nil
}

// convertTags turns Tags into prometheus labels:
func (r *Reporter) convertTags(tags metrics.Tags) prometheus.Labels {
	labels := prometheus.Labels{}
	for key, value := range tags {
		labels[key] = r.stripUnsupportedCharacters(value)
	}
	return labels
}

// listTagKeys returns a list of tag keys (we need to provide this to the Prometheus client):
func (r *Reporter) listTagKeys(tags metrics.Tags) (labelKeys []string) {
	for key := range tags {
		labelKeys = append(labelKeys, key)
	}
	return
}

// stripUnsupportedCharacters cleans up a metrics key or value:
func (r *Reporter) stripUnsupportedCharacters(metricName string) string {
	valueWithoutDots := strings.Replace(metricName, ".", "_", -1)
	valueWithoutCommas := strings.Replace(valueWithoutDots, ",", "_", -1)
	valueWithoutSpaces := strings.Replace(valueWithoutCommas, " ", "", -1)
	return valueWithoutSpaces
}
