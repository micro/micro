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
// Original source: github.com/micro/go-plugins/v3/metrics/prometheus/metric_family.go

package prometheus

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// metricFamily stores our cached metrics:
type metricFamily struct {
	counters           map[string]*prometheus.CounterVec
	gauges             map[string]*prometheus.GaugeVec
	timings            map[string]*prometheus.SummaryVec
	defaultLabels      prometheus.Labels
	mutex              sync.Mutex
	prometheusRegistry *prometheus.Registry
	timingObjectives   map[float64]float64
}

// newMetricFamily returns a new metricFamily (useful in case we want to change the structure later):
func (r *Reporter) newMetricFamily() metricFamily {

	// Take quantile thresholds from our pre-defined list:
	timingObjectives := make(map[float64]float64)
	for _, percentile := range r.options.Percentiles {
		if quantileThreshold, ok := quantileThresholds[percentile]; ok {
			timingObjectives[percentile] = quantileThreshold
		}
	}

	return metricFamily{
		counters:           make(map[string]*prometheus.CounterVec),
		gauges:             make(map[string]*prometheus.GaugeVec),
		timings:            make(map[string]*prometheus.SummaryVec),
		defaultLabels:      r.convertTags(r.options.DefaultTags),
		prometheusRegistry: r.prometheusRegistry,
		timingObjectives:   timingObjectives,
	}
}

// getCounter either gets a counter, or makes a new one:
func (mf *metricFamily) getCounter(name string, labelNames []string) *prometheus.CounterVec {
	mf.mutex.Lock()
	defer mf.mutex.Unlock()

	// See if we already have this counter:
	counter, ok := mf.counters[name]
	if !ok {

		// Make a new counter:
		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
			},
			labelNames,
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(counter)
		mf.counters[name] = counter
	}

	return counter
}

// getGauge either gets a gauge, or makes a new one:
func (mf *metricFamily) getGauge(name string, labelNames []string) *prometheus.GaugeVec {
	mf.mutex.Lock()
	defer mf.mutex.Unlock()

	// See if we already have this gauge:
	gauge, ok := mf.gauges[name]
	if !ok {

		// Make a new gauge:
		gauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
			},
			labelNames,
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(gauge)
		mf.gauges[name] = gauge
	}

	return gauge
}

// getTiming either gets a timing, or makes a new one:
func (mf *metricFamily) getTiming(name string, labelNames []string) *prometheus.SummaryVec {
	mf.mutex.Lock()
	defer mf.mutex.Unlock()

	// See if we already have this timing:
	timing, ok := mf.timings[name]
	if !ok {

		// Make a new timing:
		timing = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
				Objectives:  mf.timingObjectives,
			},
			labelNames,
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(timing)
		mf.timings[name] = timing
	}

	return timing
}
