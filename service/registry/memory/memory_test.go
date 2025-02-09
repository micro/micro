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
// Original source: github.com/micro/go-micro/v3/registry/memory/memory_test.go

package memory

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/micro/micro/v5/service/registry"
)

var (
	testData = map[string][]*registry.Service{
		"foo": {
			{
				Name:    "foo",
				Version: "1.0.0",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.0-123",
						Address: "localhost:9999",
					},
					{
						Id:      "foo-1.0.0-321",
						Address: "localhost:9999",
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.1",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.1-321",
						Address: "localhost:6666",
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.3",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.3-345",
						Address: "localhost:8888",
					},
				},
			},
		},
		"bar": {
			{
				Name:    "bar",
				Version: "default",
				Nodes: []*registry.Node{
					{
						Id:      "bar-1.0.0-123",
						Address: "localhost:9999",
					},
					{
						Id:      "bar-1.0.0-321",
						Address: "localhost:9999",
					},
				},
			},
			{
				Name:    "bar",
				Version: "latest",
				Nodes: []*registry.Node{
					{
						Id:      "bar-1.0.1-321",
						Address: "localhost:6666",
					},
				},
			},
		},
	}
)

func TestMemoryRegistry(t *testing.T) {
	m := NewRegistry()

	fn := func(k string, v []*registry.Service) {
		services, err := m.GetService(k)
		if err != nil {
			t.Errorf("Unexpected error getting service %s: %v", k, err)
		}

		if len(services) != len(v) {
			t.Errorf("Expected %d services for %s, got %d", len(v), k, len(services))
		}

		for _, service := range v {
			var seen bool
			for _, s := range services {
				if s.Version == service.Version {
					seen = true
					break
				}
			}
			if !seen {
				t.Errorf("expected to find version %s", service.Version)
			}
		}
	}

	// register data
	for _, v := range testData {
		serviceCount := 0
		for _, service := range v {
			if err := m.Register(service); err != nil {
				t.Errorf("Unexpected register error: %v", err)
			}
			serviceCount++
			// after the service has been registered we should be able to query it
			services, err := m.GetService(service.Name)
			if err != nil {
				t.Errorf("Unexpected error getting service %s: %v", service.Name, err)
			}
			if len(services) != serviceCount {
				t.Errorf("Expected %d services for %s, got %d", serviceCount, service.Name, len(services))
			}
		}
	}

	// using test data
	for k, v := range testData {
		fn(k, v)
	}

	services, err := m.ListServices()
	if err != nil {
		t.Errorf("Unexpected error when listing services: %v", err)
	}

	totalServiceCount := 0
	for _, testSvc := range testData {
		for range testSvc {
			totalServiceCount++
		}
	}

	if len(services) != totalServiceCount {
		t.Errorf("Expected total service count: %d, got: %d", totalServiceCount, len(services))
	}

	// deregister
	for _, v := range testData {
		for _, service := range v {
			if err := m.Deregister(service); err != nil {
				t.Errorf("Unexpected deregister error: %v", err)
			}
		}
	}

	// after all the service nodes have been deregistered we should not get any results
	for _, v := range testData {
		for _, service := range v {
			services, err := m.GetService(service.Name)
			if err != registry.ErrNotFound {
				t.Errorf("Expected error: %v, got: %v", registry.ErrNotFound, err)
			}
			if len(services) != 0 {
				t.Errorf("Expected %d services for %s, got %d", 0, service.Name, len(services))
			}
		}
	}
}

func TestMemoryRegistryTTL(t *testing.T) {
	m := NewRegistry()

	for _, v := range testData {
		for _, service := range v {
			if err := m.Register(service, registry.RegisterTTL(time.Millisecond)); err != nil {
				t.Fatal(err)
			}
		}
	}

	time.Sleep(ttlPruneTime * 2)

	for name := range testData {
		svcs, err := m.GetService(name)
		if err != registry.ErrNotFound {
			t.Fatal(err)
		}

		for _, svc := range svcs {
			if len(svc.Nodes) > 0 {
				t.Fatalf("Service %q still has nodes registered", name)
			}
		}
	}
}

func TestMemoryRegistryTTLConcurrent(t *testing.T) {
	concurrency := 1000
	waitTime := ttlPruneTime * 2
	m := NewRegistry()

	for _, v := range testData {
		for _, service := range v {
			if err := m.Register(service, registry.RegisterTTL(waitTime/2)); err != nil {
				t.Fatal(err)
			}
		}
	}

	if len(os.Getenv("IN_TRAVIS_CI")) == 0 {
		t.Logf("test will wait %v, then check TTL timeouts", waitTime)
	}

	errChan := make(chan error, concurrency)
	syncChan := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		go func() {
			<-syncChan
			for name := range testData {
				svcs, err := m.GetService(name)
				if err != registry.ErrNotFound {
					errChan <- err
					return
				}

				for _, svc := range svcs {
					if len(svc.Nodes) > 0 {
						errChan <- fmt.Errorf("Service %q still has nodes registered", name)
						return
					}
				}
			}

			errChan <- nil
		}()
	}

	time.Sleep(waitTime)
	close(syncChan)

	for i := 0; i < concurrency; i++ {
		if err := <-errChan; err != nil {
			t.Fatal(err)
		}
	}
}

func TestMemoryWildcard(t *testing.T) {
	m := NewRegistry()
	testSrv := &registry.Service{Name: "foo", Version: "1.0.0"}

	if err := m.Register(testSrv, registry.RegisterDomain("one")); err != nil {
		t.Fatalf("Register err: %v", err)
	}
	if err := m.Register(testSrv, registry.RegisterDomain("two")); err != nil {
		t.Fatalf("Register err: %v", err)
	}

	if recs, err := m.ListServices(registry.ListDomain("one")); err != nil {
		t.Errorf("List err: %v", err)
	} else if len(recs) != 1 {
		t.Errorf("Expected 1 record, got %v", len(recs))
	}

	if recs, err := m.ListServices(registry.ListDomain("*")); err != nil {
		t.Errorf("List err: %v", err)
	} else if len(recs) != 2 {
		t.Errorf("Expected 2 records, got %v", len(recs))
	}

	if recs, err := m.GetService(testSrv.Name, registry.GetDomain("one")); err != nil {
		t.Errorf("Get err: %v", err)
	} else if len(recs) != 1 {
		t.Errorf("Expected 1 record, got %v", len(recs))
	}

	if recs, err := m.GetService(testSrv.Name, registry.GetDomain("*")); err != nil {
		t.Errorf("Get err: %v", err)
	} else if len(recs) != 2 {
		t.Errorf("Expected 2 records, got %v", len(recs))
	}
}
