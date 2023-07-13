package registry

import (
	"net/http"

	"micro.dev/v4/service/registry"
)

// default resolver for legacy purposes
// it uses proxy routing to resolve names
// /foo becomes namespace.foo
// /v1/foo becomes namespace.v1.foo
type apiResolver struct{}

func (r *apiResolver) Resolve(req *http.Request) *Endpoint {
	// get route
	service, endpoint := apiRoute(req.URL.Path)

	// check for the namespace in the request header, this can be set by the client or injected
	// by the auth wrapper if an auth token was provided. The headr takes priority over any domain
	// passed as a default
	domain := registry.DefaultDomain

	if dom := req.Header.Get("Micro-Namespace"); len(dom) > 0 {
		domain = dom
	}

	return &Endpoint{
		Name:   service,
		Method: endpoint,
		Domain: domain,
	}
}

type Endpoint struct {
	Name   string
	Method string
	Domain string
}

func NewResolver() *apiResolver {
	return new(apiResolver)
}
