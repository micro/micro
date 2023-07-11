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
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"micro.dev/v4/service/api"
	"micro.dev/v4/service/api/router"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/registry"
	"micro.dev/v4/service/registry/cache"
	"micro.dev/v4/util/namespace"
	util "micro.dev/v4/util/router"
)

var (
	errEmptyNamespace = errors.New("namespace is empty")
	errNotFound       = errors.New("not found")
)

// endpoint struct, that holds compiled pcre
type endpoint struct {
	hostregs []*regexp.Regexp
	pathregs []util.Pattern
	pcreregs []*regexp.Regexp
}

// namespaceEntry holds the services and endpoint regexs for a namespace
type namespaceEntry struct {
	sync.RWMutex
	eps map[string]*api.Service
}

// router is the default router
type registryRouter struct {
	exit chan bool
	opts router.Options

	// registry cache
	rc cache.Cache

	// refresh channel
	refreshChan chan string

	sync.RWMutex
	namespaces map[string]*namespaceEntry

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

// refreshNamespace refreshes the list of api services in the given namespace
func (r *registryRouter) refreshNamespace(ns string) error {
	services, err := r.opts.Registry.ListServices(registry.ListDomain(ns))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("unable to list services: %v", err)
		}
		return err
	}
	if len(services) == 0 {
		return errEmptyNamespace
	}

	// for each service, get service and store endpoints
	for _, s := range services {
		// if we have nodes then use them
		dns := getDomain(s)
		if len(s.Nodes) > 0 && len(dns) > 0 {
			r.store(dns, []*registry.Service{s})
			continue
		}

		service, err := r.rc.GetService(s.Name, registry.GetDomain(ns))
		if err != nil {
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Errorf("unable to get service: %v", err)
			}
			continue
		}

		// store each independently as if we have a wildcard query
		// the domain for each service may differ
		for _, srv := range service {
			// get the namespace from the service
			ns = getDomain(srv)
			r.store(ns, []*registry.Service{srv})
		}
	}

	return nil
}

// refresh list of api services
func (r *registryRouter) refresh() {
	refreshed := make(map[string]time.Time)

	// do first load
	r.refreshNamespace(registry.WildcardDomain)

	for {
		r.RLock()
		namespaces := r.namespaces
		r.RUnlock()

		for ns, _ := range namespaces {
			err := r.refreshNamespace(ns)
			if err == errEmptyNamespace {
				r.Lock()
				delete(namespaces, ns)
				r.Unlock()
			}
		}

		// refresh the list every minute
		// TODO: rely solely on watcher
		select {
		case domain := <-r.refreshChan:
			v, ok := refreshed[domain]
			if ok && time.Since(v) < time.Minute {
				break
			}
			r.refreshNamespace(domain)
		case <-time.After(time.Minute):
		case <-r.exit:
			return
		}
	}
}

// process watch event
func (r *registryRouter) process(res *registry.Result) {
	// skip these things
	if res == nil || res.Service == nil {
		return
	}

	// get entry from cache
	// only deals with default namespace
	service, err := r.rc.GetService(res.Service.Name)
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("unable to get %v service: %v", res.Service.Name, err)
		}
		return
	}

	// only process if there's data
	if len(service) == 0 {
		return
	}

	// get te namespace
	namespace := getDomain(service[0])

	// update our local endpoints
	r.store(namespace, service)
}

// store local endpoint cache
func (r *registryRouter) store(namespace string, services []*registry.Service) {
	// endpoints
	eps := map[string]*api.Service{}

	// services
	names := map[string]bool{}

	// create a new endpoint mapping
	for _, service := range services {
		// set names we need later
		names[service.Name] = true

		// map per endpoint
		for _, sep := range service.Endpoints {
			// create a key service:endpoint_name
			key := fmt.Sprintf("%s.%s", service.Name, sep.Name)
			// decode endpoint
			end := api.Decode(sep.Metadata)
			// no endpoint or no name
			if end == nil || len(end.Name) == 0 {
				continue
			}
			// if we got nothing skip
			if err := api.Validate(end); err != nil {
				if logger.V(logger.TraceLevel, logger.DefaultLogger) {
					logger.Tracef("endpoint validation failed: %v", err)
				}
				continue
			}

			// try get endpoint
			ep, ok := eps[key]
			if !ok {
				ep = &api.Service{Name: service.Name, Domain: namespace}
			}

			// overwrite the endpoint
			ep.Endpoint = end
			// append services
			ep.Services = append(ep.Services, service)
			// store it
			eps[key] = ep
		}
	}

	r.Lock()
	nse, ok := r.namespaces[namespace]
	if !ok {
		nse = &namespaceEntry{
			eps: map[string]*api.Service{},
		}
		r.namespaces[namespace] = nse
	}
	r.Unlock()

	nse.Lock()
	defer nse.Unlock()

	// delete any existing eps for services we know
	for key, service := range nse.eps {
		// skip what we don't care about
		if !names[service.Name] {
			continue
		}

		// ok we know this thing
		// delete delete delete
		delete(nse.eps, key)
	}

	// set eps
	for name, ep := range eps {
		nse.eps[name] = ep
	}

	// save it
	r.namespaces[namespace] = nse
}

// watch for endpoint changes
func (r *registryRouter) watch() {
	var attempts int

	for {
		if r.isClosed() {
			return
		}

		// watch for changes
		w, err := r.opts.Registry.Watch(registry.WatchDomain(registry.WildcardDomain))
		if err != nil {
			attempts++
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Errorf("error watching endpoints: %v", err)
			}
			time.Sleep(time.Duration(attempts) * time.Second)
			continue
		}

		ch := make(chan bool)

		go func() {
			select {
			case <-ch:
				w.Stop()
			case <-r.exit:
				w.Stop()
			}
		}()

		// reset if we get here
		attempts = 0

		for {
			// process next event
			res, err := w.Next()
			if err != nil {
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Errorf("error getting next endpoint: %v", err)
				}
				close(ch)
				break
			}
			r.process(res)
		}
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

	// check the endpoint cache
	key := fmt.Sprintf("%s.%s", rp.Name, rp.Method)

	r.Lock()
	namespace, ok := r.namespaces[key]
	r.Unlock()

	// got the namespace
	if ok {
		if srv, kk := namespace.eps[key]; kk {
			return srv, nil
		}

		// not ok
	}

	// trigger an endpoint refresh
	select {
	case r.refreshChan <- rp.Domain:
	default:
	}

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
		exit:        make(chan bool),
		refreshChan: make(chan string),
		opts:        options,
		resolver:    new(apiResolver),
		rc:          cache.New(options.Registry),
		namespaces: map[string]*namespaceEntry{
			namespace.DefaultNamespace: &namespaceEntry{
				eps: make(map[string]*api.Service),
			}},
	}
	go r.watch()
	go r.refresh()
	return r
}

// NewRouter returns the default router
func NewRouter(opts ...router.Option) router.Router {
	return newRouter(opts...)
}
