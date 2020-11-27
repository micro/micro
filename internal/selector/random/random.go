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
// Original source: github.com/micro/go-micro/v3/selector/random/random.go

package random

import (
	"math/rand"

	"github.com/micro/micro/v3/internal/selector"
)

type random struct{}

func (r *random) Select(routes []string, opts ...selector.SelectOption) (selector.Next, error) {
	// we can't select from an empty pool of routes
	if len(routes) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	// return the next func
	return func() string {
		// if there is only one route provided we'll select it
		if len(routes) == 1 {
			return routes[0]
		}

		// select a random route from the slice
		return routes[rand.Intn(len(routes))]
	}, nil
}

func (r *random) Record(addr string, err error) error {
	return nil
}

func (r *random) Reset() error {
	return nil
}

func (r *random) String() string {
	return "random"
}

// NewSelector returns a random selector
func NewSelector(opts ...selector.Option) selector.Selector {
	return new(random)
}
