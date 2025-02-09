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
// Original source: github.com/micro/go-micro/v3/registry/memory/util.go

package memory

import (
	"time"

	"github.com/micro/micro/v5/service/registry"
)

func serviceToRecord(s *registry.Service, ttl time.Duration) *record {
	metadata := make(map[string]string, len(s.Metadata))
	for k, v := range s.Metadata {
		metadata[k] = v
	}

	nodes := make(map[string]*node, len(s.Nodes))
	for _, n := range s.Nodes {
		nodes[n.Id] = &node{
			Node:     n,
			TTL:      ttl,
			LastSeen: time.Now(),
		}
	}

	endpoints := make([]*registry.Endpoint, len(s.Endpoints))
	for i, e := range s.Endpoints {
		endpoints[i] = e
	}

	return &record{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  metadata,
		Nodes:     nodes,
		Endpoints: endpoints,
	}
}

func recordToService(r *record, domain string) *registry.Service {
	metadata := make(map[string]string, len(r.Metadata))
	for k, v := range r.Metadata {
		metadata[k] = v
	}

	// set the domain in metadata so it can be determined when a wildcard query is performed
	metadata["domain"] = domain

	endpoints := make([]*registry.Endpoint, len(r.Endpoints))
	for i, e := range r.Endpoints {
		request := new(registry.Value)
		if e.Request != nil {
			*request = *e.Request
		}
		response := new(registry.Value)
		if e.Response != nil {
			*response = *e.Response
		}

		metadata := make(map[string]string, len(e.Metadata))
		for k, v := range e.Metadata {
			metadata[k] = v
		}

		endpoints[i] = &registry.Endpoint{
			Name:     e.Name,
			Request:  request,
			Response: response,
			Metadata: metadata,
		}
	}

	nodes := make([]*registry.Node, len(r.Nodes))
	i := 0
	for _, n := range r.Nodes {
		metadata := make(map[string]string, len(n.Metadata))
		for k, v := range n.Metadata {
			metadata[k] = v
		}

		nodes[i] = &registry.Node{
			Id:       n.Id,
			Address:  n.Address,
			Metadata: metadata,
		}
		i++
	}

	return &registry.Service{
		Name:      r.Name,
		Version:   r.Version,
		Metadata:  metadata,
		Endpoints: endpoints,
		Nodes:     nodes,
	}
}
