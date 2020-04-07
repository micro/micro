package web

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/memory"
)

func TestWebResolver(t *testing.T) {
	r := memory.NewRegistry()

	selector := selector.NewSelector(
		selector.Registry(r),
	)

	res := &resolver{
		Namespace: "go.micro.web",
		Selector:  selector,
	}

	testCases := []struct {
		Host    string
		Path    string
		Service string
		Type    string
	}{
		{"localhost:8082", "/foobar", "go.micro.web.foobar", "path"},
		{"web.micro.mu", "/foobar", "go.micro.web.foobar", "path"},
		{"127.0.0.1:8082", "/hello", "go.micro.web.hello", "path"},
		{"foo.micro.mu", "/", "go.micro.web.foo", "domain"},
		{"bar.micro.mu", "/", "go.micro.web.bar", "domain"},
		{"man.web.micro.mu", "/", "go.micro.web.man", "domain"},
	}

	for _, service := range testCases {
		// set resolver type
		res.Type = service.Type

		v := &registry.Service{
			Name:    service.Service,
			Version: "latest",
			Nodes: []*registry.Node{
				{Id: "1", Address: "127.0.0.1:8080"},
			},
		}

		r.Register(v)

		u, err := url.Parse("https://" + service.Host + service.Path)
		if err != nil {
			t.Fatal(err)
		}

		req := &http.Request{
			Header: make(http.Header),
			URL:    u,
			Host:   u.Hostname(),
		}
		if endpoint, err := res.Resolve(req); err != nil {
			t.Fatalf("Failed to resolve %v: %v", service, err)
		} else if endpoint.Host != "127.0.0.1:8080" {
			t.Fatalf("Failed to resolve %v", service.Host)
		}

		r.Deregister(v)
	}

}
