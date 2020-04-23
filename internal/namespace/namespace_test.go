package namespace

import (
	"testing"

	"github.com/micro/go-micro/v2/errors"
)

func TestFromService(t *testing.T) {
	tt := []struct {
		Name      string
		Namespace string
		Error     string
	}{
		{Name: "go.micro.runtime", Namespace: RuntimeNamespace},
		{Name: "go.micro.service.foo", Namespace: "go.micro"},
		{Name: "go.micro.web.foo", Namespace: "go.micro"},
		{Name: "foo.bar", Namespace: "", Error: "Missing service type in name"},
		{Name: "service.foo", Namespace: DefaultNamespace},
		{Name: "foo.service.bar", Namespace: "foo"},
		{Name: "foo-service-bar", Namespace: "foo"},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ns, err := FromService(tc.Name)
			if ns != tc.Namespace {
				t.Errorf("Expected '%v' namespace, got '%v'", tc.Namespace, ns)
			}
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("Expected '%v' error, got nil", tc.Error)
			}
			if err == nil {
				return
			}
			mErr, ok := err.(*errors.Error)
			if !ok {
				t.Fatalf("Invalid error returned: '%v'", err)
			}
			if mErr.Detail != tc.Error {
				t.Errorf("Expected '%v' error, got '%v'", tc.Error, mErr.Detail)
			}
		})
	}
}
