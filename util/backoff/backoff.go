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
// Original source: github.com/micro/go-micro/v3/util/backoff/backoff.go

// Package backoff provides backoff functionality
package backoff

import (
	"math"
	"time"
)

// Do is a function x^e multiplied by a factor of 0.1 second.
// Result is limited to 2 minute.
func Do(attempts int) time.Duration {
	if attempts > 13 {
		return 2 * time.Minute
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100
}
