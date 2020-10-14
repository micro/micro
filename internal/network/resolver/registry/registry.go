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
// Original source: github.com/micro/go-micro/v3/network/resolver/registry/registry.go

// Package registry resolves names using the go-micro registry
package registry

import (
	"github.com/micro/micro/v3/internal/network/resolver"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/registry/mdns"
)

// Resolver is a registry network resolver
type Resolver struct {
	// Registry is the registry to use otherwise we use the defaul
	Registry registry.Registry
}

// Resolve assumes ID is a domain name e.g micro.mu
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	reg := r.Registry
	if reg == nil {
		reg = mdns.NewRegistry()
	}

	services, err := reg.GetService(name)
	if err != nil {
		return nil, err
	}

	var records []*resolver.Record

	for _, service := range services {
		for _, node := range service.Nodes {
			records = append(records, &resolver.Record{
				Address: node.Address,
			})
		}
	}

	return records, nil
}
