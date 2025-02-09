package registry

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/micro/micro/v5/service/api"
	"github.com/micro/micro/v5/service/api/router"
	"github.com/micro/micro/v5/service/context"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/registry"
	"github.com/micro/micro/v5/service/registry/cache"
	"github.com/micro/micro/v5/util/namespace"
	util "github.com/micro/micro/v5/util/router"
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
	// compiled regexp for host and path
	ceps map[string]*endpoint
}

// router is the default router
type registryRouter struct {
	exit chan bool
	opts router.Options

	// registry cache
	rc cache.Cache

	// refresh channel
	refreshChan chan string
	resolver    *apiResolver

	sync.RWMutex
	namespaces map[string]*namespaceEntry
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

		p := ep.Endpoint.Path
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
	rp := r.resolver.Resolve(req)
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
			md, ok := context.FromContext(ctx)
			if !ok {
				md = make(context.Metadata)
			}
			for k, v := range matches {
				md[fmt.Sprintf("x-api-field-%s", k)] = v
			}
			md["x-api-body"] = ep.Body
			*req = *req.Clone(context.NewContext(ctx, md))
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
	rp := r.resolver.Resolve(req)
	// service name
	name := rp.Name

	// trigger an endpoint refresh
	select {
	case r.refreshChan <- rp.Domain:
	default:
	}

	// get service
	services, err := r.rc.GetService(name, registry.GetDomain(rp.Domain))
	if err != nil {
		return nil, err
	}

	// construct api service
	return &api.Service{
		Name: name,
		Endpoint: &api.Endpoint{
			Name:    rp.Method,
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
		rc:          cache.New(options.Registry),
		resolver:    new(apiResolver),
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
