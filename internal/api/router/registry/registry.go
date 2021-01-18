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
	"strings"
	"sync"
	"time"

	"github.com/micro/micro/v3/internal/api/router"
	"github.com/micro/micro/v3/internal/namespace"
	util "github.com/micro/micro/v3/internal/router"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/registry/cache"
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
	eps map[string]*api.Service
	// compiled regexp for host and path
	ceps map[string]*endpoint
	sync.RWMutex
}

// router is the default router
type registryRouter struct {
	exit chan bool
	opts router.Options

	// registry cache
	rc cache.Cache

	sync.RWMutex
	namespaces map[string]*namespaceEntry
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
		service, err := r.rc.GetService(s.Name, registry.GetDomain(ns))
		if err != nil {
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Errorf("unable to get service: %v", err)
			}
			continue
		}
		r.store(ns, service)
	}
	return nil
}

// refresh list of api services
func (r *registryRouter) refresh() {
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
		// refresh list in 5 minutes... cruft
		// use registry watching
		select {
		case <-time.After(time.Minute * 5):
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

	// update our local endpoints
	// TODO, we currently only watch for registry changes in the default namespace
	r.store(namespace.DefaultNamespace, service)
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
				ep = &api.Service{Name: service.Name}
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
			eps:  map[string]*api.Service{},
			ceps: map[string]*endpoint{},
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

	// now set the eps we have
	for name, ep := range eps {
		nse.eps[name] = ep
		cep := &endpoint{}

		for _, h := range ep.Endpoint.Host {
			if h == "" || h == "*" {
				continue
			}
			hostreg, err := regexp.CompilePOSIX(h)
			if err != nil {
				if logger.V(logger.TraceLevel, logger.DefaultLogger) {
					logger.Tracef("endpoint have invalid host regexp: %v", err)
				}
				continue
			}
			cep.hostregs = append(cep.hostregs, hostreg)
		}

		for _, p := range ep.Endpoint.Path {
			var pcreok bool

			if p[0] == '^' && p[len(p)-1] == '$' {
				pcrereg, err := regexp.CompilePOSIX(p)
				if err == nil {
					cep.pcreregs = append(cep.pcreregs, pcrereg)
					pcreok = true
				}
			}

			rule, err := util.Parse(p)
			if err != nil && !pcreok {
				if logger.V(logger.TraceLevel, logger.DefaultLogger) {
					logger.Tracef("endpoint have invalid path pattern: %v", err)
				}
				continue
			} else if err != nil && pcreok {
				continue
			}

			tpl := rule.Compile()
			pathreg, err := util.NewPattern(tpl.Version, tpl.OpCodes, tpl.Pool, "")
			if err != nil {
				if logger.V(logger.TraceLevel, logger.DefaultLogger) {
					logger.Tracef("endpoint have invalid path pattern: %v", err)
				}
				continue
			}
			cep.pathregs = append(cep.pathregs, pathreg)
		}

		nse.ceps[name] = cep
	}
}

// watch for endpoint changes
func (r *registryRouter) watch() {
	var attempts int

	for {
		if r.isClosed() {
			return
		}

		// watch for changes
		w, err := r.opts.Registry.Watch()
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

func (r *registryRouter) Register(ep *api.Endpoint) error {
	return nil
}

func (r *registryRouter) Deregister(ep *api.Endpoint) error {
	return nil
}

func (r *registryRouter) Endpoint(req *http.Request) (*api.Service, error) {
	if r.isClosed() {
		return nil, errors.New("router closed")
	}

	var idx int
	if len(req.URL.Path) > 0 && req.URL.Path != "/" {
		idx = 1
	}
	path := strings.Split(req.URL.Path[idx:], "/")

	// resolve so we can get the namespace
	rp, err := r.opts.Resolver.Resolve(req)
	if err != nil {
		return nil, err
	}
	var ret *api.Service
	r.RLock()
	nse, ok := r.namespaces[rp.Domain]
	r.RUnlock()
	if !ok {
		// no entry in cache
		// TODO should we refresh the cache here?
		return nil, errNotFound
	}
	nse.RLock()
	defer nse.RUnlock()
endpointLoop:
	// loop through all endpoints to find either a path match or a regex match
	// prefer path matches over regexp matches e.g. prefer /foobar over ^/.*$
	// TODO: weighted matching
	for n, e := range nse.eps {
		cep, ok := nse.ceps[n]
		if !ok {
			continue
		}
		ep := e.Endpoint
		var mMatch, hMatch bool
		// 1. try method
		for _, m := range ep.Method {
			if m == req.Method {
				mMatch = true
				break
			}
		}
		if !mMatch {
			continue
		}
		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("api method match %s", req.Method)
		}

		// 2. try host
		if len(ep.Host) == 0 {
			hMatch = true
		} else {
			for idx, h := range ep.Host {
				if h == "" || h == "*" {
					hMatch = true
					break
				} else {
					if cep.hostregs[idx].MatchString(req.URL.Host) {
						hMatch = true
						break
					}
				}
			}
		}
		if !hMatch {
			continue
		}
		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("api host match %s", req.URL.Host)
		}

		// 3. try path via google.api path matching
		for _, pathreg := range cep.pathregs {
			matches, err := pathreg.Match(path, "")
			if err != nil {
				if logger.V(logger.DebugLevel, logger.DefaultLogger) {
					logger.Debugf("api gpath not match %s != %v", path, pathreg)
				}
				continue
			}
			if logger.V(logger.DebugLevel, logger.DefaultLogger) {
				logger.Debugf("api gpath match %s = %v", path, pathreg)
			}
			ctx := req.Context()
			md, ok := metadata.FromContext(ctx)
			if !ok {
				md = make(metadata.Metadata)
			}
			for k, v := range matches {
				md[fmt.Sprintf("x-api-field-%s", k)] = v
			}
			md["x-api-body"] = ep.Body
			*req = *req.Clone(metadata.NewContext(ctx, md))
			ret = e
			break endpointLoop
		}

		// 4. try path via pcre path matching
		for _, pathreg := range cep.pcreregs {
			if !pathreg.MatchString(req.URL.Path) {
				if logger.V(logger.DebugLevel, logger.DefaultLogger) {
					logger.Debugf("api pcre path not match %s != %v", path, pathreg)
				}
				continue
			}
			if logger.V(logger.DebugLevel, logger.DefaultLogger) {
				logger.Debugf("api pcre path match %s != %v", path, pathreg)
			}
			ret = e
			break
		}

		// TODO: Percentage traffic
	}
	if ret != nil {
		return ret, nil
	}

	// no match
	return nil, errNotFound
}

func (r *registryRouter) Route(req *http.Request) (*api.Service, error) {
	if r.isClosed() {
		return nil, errors.New("router closed")
	}

	// try get an endpoint from cache
	ep, err := r.Endpoint(req)
	if err == nil {
		return ep, nil
	}

	// error not nil
	// ignore that shit
	// TODO: don't ignore that shit

	// get the service name
	rp, err := r.opts.Resolver.Resolve(req)
	if err != nil {
		return nil, err
	}
	// service name
	name := rp.Name

	r.refreshNamespace(rp.Domain)
	// try to find a matching endpoint again
	if ep, err := r.Endpoint(req); err == nil {
		return ep, nil
	}

	// get service
	services, err := r.rc.GetService(name, registry.GetDomain(rp.Domain))
	if err != nil {
		return nil, err
	}

	// only use endpoint matching when the meta handler is set aka api.Default
	switch r.opts.Handler {
	// rpc handlers
	case "meta", "api", "rpc":
		handler := r.opts.Handler

		// set default handler to api
		if r.opts.Handler == "meta" {
			handler = "rpc"
		}

		// construct api service
		return &api.Service{
			Name: name,
			Endpoint: &api.Endpoint{
				Name:    rp.Method,
				Handler: handler,
			},
			Services: services,
		}, nil
	// http handler
	case "http", "proxy", "web":
		// construct api service
		return &api.Service{
			Name: name,
			Endpoint: &api.Endpoint{
				Name:    req.URL.String(),
				Handler: r.opts.Handler,
				Host:    []string{req.Host},
				Method:  []string{req.Method},
				Path:    []string{req.URL.Path},
			},
			Services: services,
		}, nil
	}

	return nil, errors.New("unknown handler")
}

func newRouter(opts ...router.Option) *registryRouter {
	options := router.NewOptions(opts...)
	r := &registryRouter{
		exit: make(chan bool),
		opts: options,
		rc:   cache.New(options.Registry),
		namespaces: map[string]*namespaceEntry{
			namespace.DefaultNamespace: &namespaceEntry{
				eps:  make(map[string]*api.Service),
				ceps: make(map[string]*endpoint),
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
