package manager

import (
	"testing"

	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/store/memory"
	"github.com/micro/micro/v3/internal/namespace"
)

func TestServices(t *testing.T) {
	testServices := []*runtime.Service{
		&runtime.Service{Name: "go.micro.service.foo", Version: "2.0.0"},
		&runtime.Service{Name: "go.micro.service.foo", Version: "1.0.0"},
		&runtime.Service{Name: "go.micro.service.bar", Version: "latest"},
	}

	testNamespace := "foo"

	m := New(&testRuntime{}, Store(memory.NewStore()), CacheStore(memory.NewStore())).(*manager)

	// listNamespaces shoud always return the default namespace
	t.Run("DefaultNamespace", func(t *testing.T) {
		if ns, err := m.listNamespaces(); err != nil {
			t.Errorf("Unexpected error when listing service: %v", err)
		} else if len(ns) != 1 {
			t.Errorf("Expected one namespace, acutually got %v", len(ns))
		} else if ns[0] != namespace.DefaultNamespace {
			t.Errorf("Expected the default namespace to be %v but was got %v", namespace.DefaultNamespace, ns[0])
		}
	})

	// creating a service should not error
	t.Run("CreateService", func(t *testing.T) {
		for _, srv := range testServices {
			if err := m.createService(srv, &runtime.CreateOptions{Namespace: testNamespace}); err != nil {
				t.Fatalf("Unexpected error when creating service %v:%v: %v", srv.Name, srv.Version, err)
			}
		}
	})

	// Calling readServices with a blank service should return all the services in the namespace
	t.Run("ReadServices", func(t *testing.T) {
		srvs, err := m.readServices(testNamespace, &runtime.Service{})
		if err != nil {
			t.Fatalf("Unexpected error when reading services%v", err)
		}
		if len(srvs) != 3 {
			t.Errorf("Expected 3 services, got %v", len(srvs))
		}
	})

	// Calling readServices with a name should return any service with that name
	t.Run("ReadServicesWithName", func(t *testing.T) {
		srvs, err := m.readServices(testNamespace, &runtime.Service{Name: "go.micro.service.foo"})
		if err != nil {
			t.Fatalf("Unexpected error when reading services%v", err)
		}
		if len(srvs) != 2 {
			t.Errorf("Expected 2 services, got %v", len(srvs))
		}
	})

	// Calling readServices with a name and version should only return the services with that name
	// and version
	t.Run("ReadServicesWithNameAndVersion", func(t *testing.T) {
		query := &runtime.Service{Name: "go.micro.service.foo", Version: "1.0.0"}
		srvs, err := m.readServices(testNamespace, query)
		if err != nil {
			t.Fatalf("Unexpected error when reading services%v", err)
		}
		if len(srvs) != 1 {
			t.Errorf("Expected 1 service, got %v", len(srvs))
		}
	})

	// Calling delete service should remove the service with that name and version
	t.Run("DeleteService", func(t *testing.T) {
		query := &runtime.Service{Name: "go.micro.service.foo", Version: "1.0.0"}
		if err := m.deleteService(testNamespace, query); err != nil {
			t.Fatalf("Unexpected error when reading services%v", err)
		}

		srvs, err := m.readServices(testNamespace, &runtime.Service{})
		if err != nil {
			t.Fatalf("Unexpected error when reading services%v", err)
		}
		if len(srvs) != 2 {
			t.Errorf("Expected 2 services, got %v", len(srvs))
		}
	})

	// a service created in one namespace shouldn't be returned when querying another
	t.Run("NamespaceScope", func(t *testing.T) {
		srv := &runtime.Service{Name: "go.micro.service.apple", Version: "latest"}

		if err := m.createService(srv, &runtime.CreateOptions{Namespace: "random"}); err != nil {
			t.Fatalf("Unexpected error when creating service %v", err)
		}

		if srvs, err := m.readServices(testNamespace, srv); err != nil {
			t.Fatalf("Unexpected error when listing services %v", err)
		} else if len(srvs) != 0 {
			t.Errorf("Expected 0 services, got %v", len(srvs))
		}
	})
}
