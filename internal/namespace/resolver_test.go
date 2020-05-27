package namespace

import (
	"net/http"
	"net/url"
	"testing"
)

func TestResolveWithType(t *testing.T) {
	tt := []struct {
		Name        string
		Namespace   string
		ServiceType string
		Host        string
		Result      string
	}{
		{
			Name:        "A fixed namespace with web type",
			ServiceType: "web",
			Namespace:   "foobar",
			Host:        "micro.mu",
			Result:      "foobar.web",
		},
		{
			Name:        "A dynamic namespace with a type web and the micro.mu domain",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "micro.mu",
			Result:      DefaultNamespace + ".web",
		},
		{
			Name:        "A dynamic namespace with a type api and the micro.mu domain",
			ServiceType: "api",
			Namespace:   "domain",
			Host:        "micro.mu",
			Result:      DefaultNamespace + ".api",
		},
		{
			Name:        "A dynamic namespace with a type web and a effective top level domain",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "micro.com.au",
			Result:      DefaultNamespace + ".web",
		},
		{
			Name:        "A dynamic namespace with a type web and the web.micro.mu subdomain",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "web.micro.mu",
			Result:      DefaultNamespace + ".web",
		},
		{
			Name:        "A dynamic namespace with a type web and a micro.mu subdomain",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "foo.micro.mu",
			Result:      DefaultNamespace + ".web",
		},
		{
			Name:        "A dynamic namespace with a type web and top level domain host",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "myapp.com",
			Result:      DefaultNamespace + ".web",
		},
		{
			Name:        "A dynamic namespace with a type web subdomain host",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "staging.myapp.com",
			Result:      "staging.web",
		},
		{
			Name:        "A dynamic namespace with a type web and multi-level subdomain host",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "staging.myapp.m3o.app",
			Result:      "myapp.staging.web",
		},
		{
			Name:        "A dynamic namespace with a type web and dev host",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "127.0.0.1",
			Result:      DefaultNamespace + ".web",
		},
		{
			Name:        "A dynamic namespace with a type web and localhost host",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "localhost",
			Result:      DefaultNamespace + ".web",
		},
		{
			Name:        "A dynamic namespace with a type web and IP host",
			ServiceType: "web",
			Namespace:   "domain",
			Host:        "81.151.101.146",
			Result:      DefaultNamespace + ".web",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			r := NewResolver(tc.ServiceType, tc.Namespace)
			result := r.ResolveWithType(&http.Request{URL: &url.URL{Host: tc.Host}})
			if result != tc.Result {
				t.Errorf("Expected namespace %v for host %v with service type %v and namespace %v, actually got %v", tc.Result, tc.Host, tc.ServiceType, tc.Namespace, result)
			}
		})
	}
}
