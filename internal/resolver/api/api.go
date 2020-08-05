// Package micro provides a micro rpc resolver which prefixes a namespace
package micro

import (
	"net/http"

	"github.com/micro/go-micro/v3/api/resolver"
)

// default resolver for legacy purposes
// it uses proxy routing to resolve names
// /foo becomes namespace.foo
// /v1/foo becomes namespace.v1.foo
type Resolver struct {
	opts resolver.Options
}

func (r *Resolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	options := resolver.NewResolveOptions(opts...)

	var name, method string

	switch r.opts.Handler {
	// internal handlers
	case "meta", "api", "rpc", "micro":
		name, method = apiRoute(req.URL.Path)
	default:
		method = req.Method
		name = proxyRoute(req.URL.Path)
	}

	// append the service prefix, e.g. foo.api
	if len(r.opts.ServicePrefix) > 0 {
		name = r.opts.ServicePrefix + "." + name
	}

	return &resolver.Endpoint{
		Name:   name,
		Domain: options.Domain,
		Method: method,
	}, nil
}

func (r *Resolver) String() string {
	return "micro"
}

// NewResolver creates a new micro resolver
func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &Resolver{
		opts: resolver.NewOptions(opts...),
	}
}
