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
// Original source: github.com/micro/micro/v3/metrics/options_tests.go

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {

	// Make some new options:
	options := NewOptions(
		Address(":9999"),
		DefaultTags(map[string]string{"service": "prometheus-test"}),
		Path("/prometheus"),
		Percentiles([]float64{0.11, 0.22, 0.33}),
	)

	// Check that the defaults and overrides were accepted:
	assert.Equal(t, ":9999", options.Address)
	assert.Equal(t, "prometheus-test", options.DefaultTags["service"])
	assert.Equal(t, "/prometheus", options.Path)
	assert.Equal(t, []float64{0.11, 0.22, 0.33}, options.Percentiles)
}
