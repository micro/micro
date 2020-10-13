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
// Original source: github.com/micro/go-micro/v3/network/resolver/static/static.go

// Package static is a static resolver
package static

import (
	"github.com/micro/micro/v3/internal/network/resolver"
)

// Resolver returns a static list of nodes. In the event the node list
// is not present it will return the name of the network passed in.
type Resolver struct {
	// A static list of nodes
	Nodes []string
}

// Resolve returns the list of nodes
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	// if there are no nodes just return the name
	if len(r.Nodes) == 0 {
		return []*resolver.Record{
			{Address: name},
		}, nil
	}

	records := make([]*resolver.Record, 0, len(r.Nodes))

	for _, node := range r.Nodes {
		records = append(records, &resolver.Record{
			Address: node,
		})
	}

	return records, nil
}
