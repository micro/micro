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
// Original source: github.com/micro/go-micro/v3/selector/tests.go

package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests runs all the tests against a selector to ensure the implementations are consistent
func Tests(t *testing.T, s Selector) {
	r1 := "127.0.0.1:8000"
	r2 := "127.0.0.1:8001"

	t.Run("Select", func(t *testing.T) {
		t.Run("NoRoutes", func(t *testing.T) {
			_, err := s.Select([]string{})
			assert.Equal(t, ErrNoneAvailable, err, "Expected error to be none available")
		})

		t.Run("OneRoute", func(t *testing.T) {
			next, err := s.Select([]string{r1})
			srv := next()
			assert.Nil(t, err, "Error should be nil")
			assert.Equal(t, r1, srv, "Expected the route to be returned")
		})

		t.Run("MultipleRoutes", func(t *testing.T) {
			next, err := s.Select([]string{r1, r2})
			assert.Nil(t, err, "Error should be nil")
			srv := next()
			if srv != r1 && srv != r2 {
				t.Errorf("Expected the route to be one of the inputs")
			}
		})
	})

	t.Run("Record", func(t *testing.T) {
		err := s.Record(r1, nil)
		assert.Nil(t, err, "Expected the error to be nil")
	})

	t.Run("String", func(t *testing.T) {
		assert.NotEmpty(t, s.String(), "String returned a blank string")
	})
}
