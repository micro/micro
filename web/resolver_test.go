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
		host    string
		path    string
		service string
	}{
		{"localhost:8082", "/foobar", "go.micro.web.foobar"},
		{"web.micro.mu", "/foobar", "go.micro.web.foobar"},
		{"127.0.0.1:8082", "/hello", "go.micro.web.hello"},
		{"foo.micro.mu", "/", "go.micro.web.foo"},
		{"bar.micro.mu", "/", "go.micro.web.bar"},
		{"man.web.micro.mu", "/", "go.micro.web.man"},
	}

	for _, service := range testCases {
		v := &registry.Service{
			Name:    service.service,
			Version: "latest",
			Nodes: []*registry.Node{
				{Id: "1", Address: "127.0.0.1:8080"},
			},
		}

		r.Register(v)

		u, err := url.Parse("https://" + service.host + service.path)
		if err != nil {
			t.Fatal(err)
		}

		req := &http.Request{
			Header: make(http.Header),
			URL:    u,
		}
		if err := res.Resolve(req); err != nil {
			t.Fatal(err)
		}

		if req.Host != "127.0.0.1:8080" {
			t.Fatalf("Failed to resolve %v", service.host)
		}

		r.Deregister(v)
	}

}
