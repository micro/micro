package web

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	// load the cmd package to load defaults since we're using a test profile without importing
	// micro or service
	_ "github.com/micro/micro/v3/cmd"

	"github.com/micro/micro/v3/internal/api/resolver"
	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/router"
)

type testCase struct {
	Host    string
	Path    string
	Service string
}

func TestWebResolver(t *testing.T) {
	profile.Test.Setup(nil)

	t.Run("WithServicePrefix", func(t *testing.T) {
		res := &Resolver{
			Options: resolver.NewOptions(
				resolver.WithServicePrefix("web"),
			),
			Router: router.DefaultRouter,
		}

		testCases := []testCase{
			{"localhost:8082", "/foobar", "web.foobar"},
			{"web.micro.mu", "/foobar", "web.foobar"},
			{"127.0.0.1:8082", "/hello", "web.hello"},
			{"demo.m3o.app", "/bar", "web.bar"},
		}

		runTests(t, res, testCases)
	})

	t.Run("WithoutServicePrefix", func(t *testing.T) {
		res := &Resolver{
			Options: resolver.NewOptions(),
			Router:  router.DefaultRouter,
		}

		testCases := []testCase{
			{"localhost:8082", "/foobar", "foobar"},
			{"web.micro.mu", "/foobar", "foobar"},
			{"127.0.0.1:8082", "/hello", "hello"},
			{"demo.m3o.app", "/bar", "bar"},
		}

		runTests(t, res, testCases)
	})
}

func runTests(t *testing.T, res *Resolver, testCases []testCase) {
	for _, service := range testCases {
		t.Run(service.Host+service.Path, func(t *testing.T) {
			v := &registry.Service{
				Name:    service.Service,
				Version: "latest",
				Nodes: []*registry.Node{
					{Id: "1", Address: "127.0.0.1:8080"},
				},
			}

			registry.DefaultRegistry.Register(v)

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

			registry.DefaultRegistry.Deregister(v)
		})
	}
}
