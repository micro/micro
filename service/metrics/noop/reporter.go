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
// Original source: github.com/micro/micro/v3/metrics/noop/reporter.go

package noop

import (
	"time"

	"github.com/micro/micro/v3/service/metrics"
)

// Reporter is an implementation of metrics.Reporter:
type Reporter struct {
	options metrics.Options
}

// New returns a configured noop reporter:
func New(opts ...metrics.Option) *Reporter {
	return &Reporter{
		options: metrics.NewOptions(opts...),
	}
}

// Count implements the metrics.Reporter interface Count method:
func (r *Reporter) Count(metricName string, value int64, tags metrics.Tags) error {
	return nil
}

// Gauge implements the metrics.Reporter interface Gauge method:
func (r *Reporter) Gauge(metricName string, value float64, tags metrics.Tags) error {
	return nil
}

// Timing implements the metrics.Reporter interface Timing method:
func (r *Reporter) Timing(metricName string, value time.Duration, tags metrics.Tags) error {
	return nil
}
