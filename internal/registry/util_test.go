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
// Original source: github.com/micro/go-micro/v3/util/registry/util_test.go

package registry

import (
	"os"
	"testing"

	"github.com/micro/micro/v3/service/registry"
)

func TestRemove(t *testing.T) {
	services := []*registry.Service{
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*registry.Node{
				{
					Id:      "foo-123",
					Address: "localhost:9999",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*registry.Node{
				{
					Id:      "foo-123",
					Address: "localhost:6666",
				},
			},
		},
	}

	servs := Remove([]*registry.Service{services[0]}, []*registry.Service{services[1]})
	if i := len(servs); i > 0 {
		t.Errorf("Expected 0 nodes, got %d: %+v", i, servs)
	}
	if len(os.Getenv("IN_TRAVIS_CI")) == 0 {
		t.Logf("Services %+v", servs)
	}
}

func TestRemoveNodes(t *testing.T) {
	services := []*registry.Service{
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*registry.Node{
				{
					Id:      "foo-123",
					Address: "localhost:9999",
				},
				{
					Id:      "foo-321",
					Address: "localhost:6666",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*registry.Node{
				{
					Id:      "foo-123",
					Address: "localhost:6666",
				},
			},
		},
	}

	nodes := delNodes(services[0].Nodes, services[1].Nodes)
	if i := len(nodes); i != 1 {
		t.Errorf("Expected only 1 node, got %d: %+v", i, nodes)
	}
	if len(os.Getenv("IN_TRAVIS_CI")) == 0 {
		t.Logf("Nodes %+v", nodes)
	}
}
