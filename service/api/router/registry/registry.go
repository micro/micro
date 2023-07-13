// Copyright 2020 Asim Aslam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/api/router/registry/registry.go

// Package registry provides a dynamic api service router
package registry

import (
	"errors"
	"net/http"

	"micro.dev/v4/service/api"
	"micro.dev/v4/service/api/router"
	"micro.dev/v4/service/registry"
	"micro.dev/v4/service/registry/cache"
)

var (
	errEmptyNamespace = errors.New("namespace is empty")
	errNotFound       = errors.New("not found")
)

// router is the default router
type registryRouter struct {
	exit chan bool
	opts router.Options

	// registry cache
	rc cache.Cache
	// api resolver
	resolver *apiResolver
}

func getDomain(srv *registry.Service) string {
	// check the service metadata for domain
	// TODO: domain as Domain field in registry?
	if srv.Metadata != nil && len(srv.Metadata["domain"]) > 0 {
		return srv.Metadata["domain"]
	} else if nodes := srv.Nodes; len(nodes) > 0 && nodes[0].Metadata != nil {
		// only return the domain if its set
		if len(nodes[0].Metadata["domain"]) > 0 {
			return nodes[0].Metadata["domain"]
		}
	}

	// otherwise return wildcard
	// TODO: return GlobalDomain or PublicDomain
	return registry.DefaultDomain
}

func (r *registryRouter) isClosed() bool {
	select {
	case <-r.exit:
		return true
	default:
		return false
	}
}

func (r *registryRouter) Options() router.Options {
	return r.opts
}

func (r *registryRouter) Close() error {
	select {
	case <-r.exit:
		return nil
	default:
		close(r.exit)
		r.rc.Stop()
	}
	return nil
}

func (r *registryRouter) Route(req *http.Request) (*api.Service, error) {
	if r.isClosed() {
		return nil, errors.New("router closed")
	}

	// get the service name
	rp := r.resolver.Resolve(req)

	// get service
	services, err := r.rc.GetService(rp.Name, registry.GetDomain(rp.Domain))
	if err != nil {
		return nil, err
	}

	// construct api service
	return &api.Service{
		Name:   rp.Name,
		Domain: rp.Domain,
		Endpoint: &api.Endpoint{
			Name:    rp.Method,
			Domain:  rp.Domain,
			Handler: "rpc",
		},
		Services: services,
	}, nil
}

func newRouter(opts ...router.Option) *registryRouter {
	options := router.NewOptions(opts...)
	r := &registryRouter{
		exit:     make(chan bool),
		opts:     options,
		resolver: new(apiResolver),
		rc:       cache.New(options.Registry),
	}
	return r
}

// NewRouter returns the default router
func NewRouter(opts ...router.Option) router.Router {
	return newRouter(opts...)
}
