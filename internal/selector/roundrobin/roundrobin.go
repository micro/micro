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
// Original source: github.com/micro/go-micro/v3/selector/roundrobin/roundrobin.go

package roundrobin

import (
	"math/rand"

	"github.com/micro/micro/v3/internal/selector"
)

// NewSelector returns an initalised round robin selector
func NewSelector(opts ...selector.Option) selector.Selector {
	return new(roundrobin)
}

type roundrobin struct{}

func (r *roundrobin) Select(routes []string, opts ...selector.SelectOption) (selector.Next, error) {
	if len(routes) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	i := rand.Intn(len(routes))

	return func() string {
		route := routes[i%len(routes)]
		// increment
		i++
		return route
	}, nil
}

func (r *roundrobin) Record(addr string, err error) error { return nil }

func (r *roundrobin) Reset() error { return nil }

func (r *roundrobin) String() string {
	return "roundrobin"
}
