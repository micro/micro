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
// Original source: github.com/micro/go-micro/v3/selector/roundrobin/roundrobin_test.go

package roundrobin

import (
	"testing"

	"github.com/micro/micro/v3/internal/selector"
	"github.com/stretchr/testify/assert"
)

func TestRoundRobin(t *testing.T) {
	selector.Tests(t, NewSelector())

	r1 := "127.0.0.1:8000"
	r2 := "127.0.0.1:8001"
	r3 := "127.0.0.1:8002"

	sel := NewSelector()

	// By passing r1 and r2 first, it forces a set sequence of (r1 => r2 => r3 => r1)

	next, err := sel.Select([]string{r1})
	r := next()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, r1, r, "Expected route to be r1")

	next, err = sel.Select([]string{r2})
	r = next()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, r2, r, "Expected route to be r2")

	routes := []string{r1, r2, r3}
	next, err = sel.Select(routes)
	assert.Nil(t, err, "Error should be nil")
	n1, n2, n3, n4 := next(), next(), next(), next()

	// start element is random but then it should loop through in order
	start := -1
	for i := 0; i < 3; i++ {
		if n1 == routes[i] {
			start = i
			break
		}
	}
	assert.NotEqual(t, start, -1)
	assert.Equal(t, routes[start], n1, "Unexpected route")
	assert.Equal(t, routes[(start+1)%3], n2, "Unexpected route")
	assert.Equal(t, routes[(start+2)%3], n3, "Unexpected route")
	assert.Equal(t, routes[(start+3)%3], n4, "Unexpected route")
}
