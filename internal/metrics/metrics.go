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
// Original source: github.com/micro/micro/v3/metrics/metrics.go

// Package metrics is for instrumentation and debugging
package metrics

import "time"

// Tags is a map of fields to add to a metric:
type Tags map[string]string

// Reporter is an interface for collecting and instrumenting metrics
type Reporter interface {
	Count(id string, value int64, tags Tags) error
	Gauge(id string, value float64, tags Tags) error
	Timing(id string, value time.Duration, tags Tags) error
}
