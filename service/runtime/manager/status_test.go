package manager

import (
	"testing"

	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/store/memory"
	"github.com/micro/micro/v2/internal/namespace"
)

type testRuntime struct {
	readServices []*runtime.Service
	runtime.Runtime
}

func (r *testRuntime) Read(...runtime.ReadOption) ([]*runtime.Service, error) {
	return r.readServices, nil
}

func TestStatus(t *testing.T) {
	testServices := []*runtime.Service{
		&runtime.Service{
			Name:     "go.micro.service.foo",
			Version:  "latest",
			Metadata: map[string]string{"status": "starting"},
		},
		&runtime.Service{
			Name:     "go.micro.service.bar",
			Version:  "2.0.0",
			Metadata: map[string]string{"status": "error", "error": "Crashed on L1"},
		},
	}

	rt := &testRuntime{readServices: testServices}
	m := New(rt, Store(memory.NewStore())).(*manager)

	// sync the status with the runtime, this should set the status for the testServices in the cache
	m.syncStatus()

	// get the statuses from the service
	statuses, err := m.listStatuses(namespace.DefaultNamespace)
	if err != nil {
		t.Fatalf("Unexpected error when listing statuses: %v", err)
	}

	// loop through the test services and check the status matches what was set in the metadata
	for _, srv := range testServices {
		s, ok := statuses[srv.Name+":"+srv.Version]
		if !ok {
			t.Errorf("Missing status for %v:%v", srv.Name, srv.Version)
			continue
		}
		if s.Status != srv.Metadata["status"] {
			t.Errorf("Incorrect status for %v:%v, expepcted %v but got %v", srv.Name, srv.Version, srv.Metadata["status"], s.Status)
		}
		if s.Error != srv.Metadata["error"] {
			t.Errorf("Incorrect error for %v:%v, expepcted %v but got %v", srv.Name, srv.Version, srv.Metadata["error"], s.Error)
		}
	}

	// update the status for a service and check it correctly updated
	srv := testServices[0]
	srv.Metadata["status"] = "running"
	if err := m.cacheStatus(namespace.DefaultNamespace, srv); err != nil {
		t.Fatalf("Unexpected error when caching status: %v", err)
	}

	// get the statuses from the service
	statuses, err = m.listStatuses(namespace.DefaultNamespace)
	if err != nil {
		t.Fatalf("Unexpected error when listing statuses: %v", err)
	}

	// check the new status matches the changed service
	s, ok := statuses[srv.Name+":"+srv.Version]
	if !ok {
		t.Errorf("Missing status for %v:%v", srv.Name, srv.Version)
	}
	if s.Status != srv.Metadata["status"] {
		t.Errorf("Incorrect status for %v:%v, expepcted %v but got %v", srv.Name, srv.Version, srv.Metadata["status"], s.Status)
	}
	if s.Error != srv.Metadata["error"] {
		t.Errorf("Incorrect error for %v:%v, expepcted %v but got %v", srv.Name, srv.Version, srv.Metadata["error"], s.Error)
	}
}
