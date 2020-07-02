package web

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/api/resolver"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/memory"
	"github.com/micro/go-micro/v2/router"
)

func TestWebResolver(t *testing.T) {
	r := memory.NewRegistry()

	res := &Resolver{
		Options: resolver.NewOptions(
			resolver.WithServicePrefix("go.micro.web"),
		),
		Router: router.NewRouter(router.Registry(r)),
	}

	testCases := []struct {
		Host    string
		Path    string
		Service string
	}{
		{"localhost:8082", "/foobar", "go.micro.web.foobar"},
		{"web.micro.mu", "/foobar", "go.micro.web.foobar"},
		{"127.0.0.1:8082", "/hello", "go.micro.web.hello"},
		{"demo.m3o.app", "/bar", "go.micro.web.bar"},
	}

	for _, service := range testCases {
		t.Run(service.Host+service.Path, func(t *testing.T) {
			v := &registry.Service{
				Name:    service.Service,
				Version: "latest",
				Nodes: []*registry.Node{
					{Id: "1", Address: "127.0.0.1:8080"},
				},
			}

			r.Register(v)

			// registry events are published to the router async (although if we don't wait the fallback should still kick in)
			time.Sleep(time.Millisecond * 10)

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
		})
	}

}
