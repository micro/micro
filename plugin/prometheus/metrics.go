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
// Original source: github.com/micro/go-plugins/v3/metrics/prometheus/metrics.go

package prometheus

import (
	"errors"
	"time"

	"github.com/micro/micro/v3/service/metrics"
)

// ErrPrometheusPanic is a catch-all for the panics which can be thrown by the Prometheus client:
var ErrPrometheusPanic = errors.New("The Prometheus client panicked. Did you do something like change the tag cardinality or the type of a metric?")

// Count is a counter with key/value tags:
// New values are added to any previous one (eg "number of hits")
func (r *Reporter) Count(name string, value int64, tags metrics.Tags) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPrometheusPanic
		}
	}()

	counter := r.metrics.getCounter(r.stripUnsupportedCharacters(name), r.listTagKeys(tags))
	metric, err := counter.GetMetricWith(r.convertTags(tags))
	if err != nil {
		return err
	}

	metric.Add(float64(value))
	return err
}

// Gauge is a register with key/value tags:
// New values simply override any previous one (eg "current connections")
func (r *Reporter) Gauge(name string, value float64, tags metrics.Tags) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPrometheusPanic
		}
	}()

	gauge := r.metrics.getGauge(r.stripUnsupportedCharacters(name), r.listTagKeys(tags))
	metric, err := gauge.GetMetricWith(r.convertTags(tags))
	if err != nil {
		return err
	}

	metric.Set(value)
	return err
}

// Timing is a histogram with key/valye tags:
// New values are added into a series of aggregations
func (r *Reporter) Timing(name string, value time.Duration, tags metrics.Tags) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPrometheusPanic
		}
	}()

	timing := r.metrics.getTiming(r.stripUnsupportedCharacters(name), r.listTagKeys(tags))
	metric, err := timing.GetMetricWith(r.convertTags(tags))
	if err != nil {
		return err
	}

	metric.Observe(value.Seconds())
	return err
}
