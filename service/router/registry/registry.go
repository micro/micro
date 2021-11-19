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
// Original source: github.com/micro/micro/v3/router/registry/registry.go

package registry

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/router"
)

var (
	// RefreshInterval is the time at which we completely refresh the table
	RefreshInterval = time.Second * 30
)

// router implements router interface
type registryRouter struct {
	sync.RWMutex

	running  bool
	table    *table
	options  router.Options
	exit     chan bool
	initChan chan bool
}

// NewRouter creates new router and returns it
func NewRouter(opts ...router.Option) router.Router {
	// get default options
	options := router.DefaultOptions()

	// apply requested options
	for _, o := range opts {
		o(&options)
	}

	// construct the router
	r := &registryRouter{
		options:  options,
		initChan: make(chan bool, 1),
		table:    newTable(),
		exit:     make(chan bool),
	}

	// initialise the router
	r.Init()

	return r
}

// Init initializes router with given options
func (r *registryRouter) Init(opts ...router.Option) error {
	r.Lock()
	for _, o := range opts {
		o(&r.options)
	}
	r.Unlock()

	// add default gateway into routing table
	if r.options.Gateway != "" {
		// note, the only non-default value is the gateway
		route := router.Route{
			Service: "*",
			Address: "*",
			Gateway: r.options.Gateway,
			Network: "*",
			Router:  r.options.Id,
			Link:    router.DefaultLink,
			Metric:  router.DefaultMetric,
		}
		if err := r.table.Create(route); err != nil && err != router.ErrDuplicateRoute {
			return fmt.Errorf("failed adding default gateway route: %s", err)
		}
	}

	// only cache if told to do so
	if r.options.Cache {
		r.run()
	}

	// push a message to the init chan so the watchers
	// can reset in the case the registry was changed
	go func() {
		select {
		case r.initChan <- true:
		default:
		}
	}()

	return nil
}

// Options returns router options
func (r *registryRouter) Options() router.Options {
	r.RLock()
	defer r.RUnlock()

	options := r.options

	return options
}

// Table returns routing table
func (r *registryRouter) Table() router.Table {
	r.Lock()
	defer r.Unlock()
	return r.table
}

func getDomain(srv *registry.Service) string {
	// check the service metadata for domain
	// TODO: domain as Domain field in registry?
	if srv.Metadata != nil && len(srv.Metadata["domain"]) > 0 {
		return srv.Metadata["domain"]
	} else if len(srv.Nodes) > 0 && srv.Nodes[0].Metadata != nil {
		return srv.Nodes[0].Metadata["domain"]
	}

	// otherwise return wildcard
	// TODO: return GlobalDomain or PublicDomain
	return registry.DefaultDomain
}

// manageRoute applies action on a given route
func (r *registryRouter) manageRoute(route router.Route, action string) error {
	switch action {
	case "create":
		if err := r.table.Create(route); err != nil && err != router.ErrDuplicateRoute {
			return fmt.Errorf("failed adding route for service %s: %s", route.Service, err)
		}
	case "delete":
		if err := r.table.Delete(route); err != nil && err != router.ErrRouteNotFound {
			return fmt.Errorf("failed deleting route for service %s: %s", route.Service, err)
		}
	case "update":
		if err := r.table.Update(route); err != nil {
			return fmt.Errorf("failed updating route for service %s: %s", route.Service, err)
		}
	default:
		return fmt.Errorf("failed to manage route for service %s: unknown action %s", route.Service, action)
	}

	return nil
}

// createRoutes turns a service into a list routes basically converting nodes to routes
func (r *registryRouter) createRoutes(service *registry.Service, network string) []router.Route {
	var routes []router.Route

	for _, node := range service.Nodes {
		routes = append(routes, router.Route{
			Service:  service.Name,
			Address:  node.Address,
			Gateway:  "",
			Network:  network,
			Router:   r.options.Id,
			Link:     router.DefaultLink,
			Metric:   router.DefaultMetric,
			Metadata: node.Metadata,
		})
	}

	return routes
}

// manageServiceRoutes applies action to all routes of the service.
// It returns error of the action fails with error.
func (r *registryRouter) manageRoutes(service *registry.Service, action, network string) error {
	// action is the routing table action
	action = strings.ToLower(action)

	// create a set of routes from the service
	routes := r.createRoutes(service, network)

	// if its a delete action and there's no nodes
	// it means we need to wipe out all the routes
	// for that service
	if action == "delete" && len(routes) == 0 {
		// delete the service entirely
		r.table.deleteService(service.Name, network)
		return nil
	}

	// create the routes in the table
	for _, route := range routes {
		logger.Tracef("Creating route %v domain: %v", route, network)
		if err := r.manageRoute(route, action); err != nil {
			return err
		}
	}

	return nil
}

// loadRoutes applies action to all routes of each service found in the registry.
// It returns error if either the services failed to be listed or the routing table action fails.
func (r *registryRouter) loadRoutes(name, domain string) error {
	var services []*registry.Service
	var err error

	if len(domain) == 0 {
		domain = registry.WildcardDomain
	}

	if len(name) > 0 {
		services, err = r.options.Registry.GetService(name, registry.GetDomain(domain))
	} else {
		services, err = r.options.Registry.ListServices(registry.ListDomain(domain))
	}

	if err != nil {
		return fmt.Errorf("failed listing services: %v", err)
	}

	// delete the services first
	for _, service := range services {
		// get the services domain from metadata. Fallback to wildcard.
		domain := getDomain(service)

		// delete the existing service
		r.table.deleteService(service.Name, domain)
	}

	// add each service version as a separate set of routes
	for _, service := range services {
		// get the services domain from metadata. Fallback to wildcard.
		domain := getDomain(service)

		// create the routes
		routes := r.createRoutes(service, domain)

		// if the routes exist save them
		if len(routes) > 0 {
			logger.Tracef("Creating routes for service %v domain: %v", service, domain)
			for _, rt := range routes {
				err := r.table.Create(rt)

				// update the route to prevent it from expiring
				if err == router.ErrDuplicateRoute {
					err = r.table.Update(rt)
				}

				if err != nil {
					logger.Errorf("Error creating route for service %v in domain %v: %v", service, domain, err)
				}
			}
			continue
		}

		// otherwise get all the service info

		// get the service to retrieve all its info
		srvs, err := r.options.Registry.GetService(service.Name, registry.GetDomain(domain))
		if err != nil {
			logger.Tracef("Failed to get service %s domain: %s", service.Name, domain)
			continue
		}

		// manage the routes for all returned services
		for _, srv := range srvs {
			routes := r.createRoutes(srv, domain)

			if len(routes) > 0 {
				logger.Tracef("Creating routes for service %v domain: %v", srv, domain)
				for _, rt := range routes {
					err := r.table.Create(rt)

					// update the route to prevent it from expiring
					if err == router.ErrDuplicateRoute {
						err = r.table.Update(rt)
					}

					if err != nil {
						logger.Errorf("Error creating route for service %v in domain %v: %v", service, domain, err)
					}
				}
			}
		}
	}

	return nil
}

// Close the router
func (r *registryRouter) Close() error {
	r.Lock()
	defer r.Unlock()

	select {
	case <-r.exit:
		return nil
	default:
		if !r.running {
			return nil
		}
		close(r.exit)

	}

	r.running = false
	return nil
}

// lookup retrieves all the routes for a given service and creates them in the routing table
func (r *registryRouter) Lookup(service string, opts ...router.LookupOption) ([]router.Route, error) {
	q := router.NewLookup(opts...)

	// if we find the routes filter and return them
	routes, err := r.table.Read(router.ReadService(service))
	if err == nil {
		routes = router.Filter(routes, q)
		if len(routes) == 0 {
			return nil, router.ErrRouteNotFound
		}
		return routes, nil
	}

	// lookup the route
	logger.Tracef("Fetching route for %s domain: %v", service, registry.WildcardDomain)

	services, err := r.options.Registry.GetService(service, registry.GetDomain(registry.WildcardDomain))
	if err == registry.ErrNotFound {
		logger.Tracef("Failed to find route for %s", service)
		return nil, router.ErrRouteNotFound
	} else if err != nil {
		logger.Tracef("Failed to find route for %s: %v", service, err)
		return nil, fmt.Errorf("failed getting services: %v", err)
	}

	for _, srv := range services {
		domain := getDomain(srv)
		// TODO: should we continue to send the event indicating we created a route?
		// lookup is only called in the query path so probably not
		routes = append(routes, r.createRoutes(srv, domain)...)
	}

	// if we're supposed to cache then save the routes
	if r.options.Cache {
		for _, route := range routes {
			r.table.Create(route)
		}
	}

	routes = router.Filter(routes, q)
	if len(routes) == 0 {
		return nil, router.ErrRouteNotFound
	}
	return routes, nil
}

// watchRegistry watches registry and updates routing table based on the received events.
// It returns error if either the registry watcher fails with error or if the routing table update fails.
func (r *registryRouter) watchRegistry(w registry.Watcher) error {
	exit := make(chan bool)

	defer func() {
		close(exit)
	}()

	go func() {
		defer w.Stop()

		select {
		case <-exit:
			return
		case <-r.initChan:
			return
		case <-r.exit:
			return
		}
	}()

	for {
		// get the next service
		res, err := w.Next()
		if err != nil {
			if err != registry.ErrWatcherStopped {
				return err
			}
			break
		}

		// don't process nil entries
		if res.Service == nil {
			logger.Trace("Received a nil service")
			continue
		}

		logger.Tracef("Router dealing with next event %s %+v\n", res.Action, res.Service)

		// we only use the registry notifications as events
		// then go on to actually query it for the full list

		// get the services domain from metadata. Fallback to wildcard.
		domain := getDomain(res.Service)

		// load routes for this service
		if err := r.loadRoutes(res.Service.Name, domain); err != nil {
			return err
		}
	}

	return nil
}

// start the router. Should be called under lock.
func (r *registryRouter) run() error {
	if r.running {
		return nil
	}

	// set running
	r.running = true

	// create a refresh notify channel
	refresh := make(chan bool, 1)

	// fires the refresh for loading routes
	refreshRoutes := func() {
		select {
		case refresh <- true:
		default:
		}
	}

	// refresh all the routes in the event of a failure watching the registry
	go func() {
		var lastRefresh time.Time

		// load a refresh
		refreshRoutes()

		for {
			select {
			case <-r.exit:
				return
			case <-refresh:
				// load new routes
				if err := r.loadRoutes("", ""); err != nil {
					logger.Debugf("failed refreshing registry routes: %s", err)
					// in this don't prune
					continue
				}

				// first time so nothing to prune
				if !lastRefresh.IsZero() {
					// prune any routes since last refresh since we've
					// updated basically everything we care about
					r.table.pruneRoutes(time.Since(lastRefresh))
				}

				// update the refresh time
				lastRefresh = time.Now()
			case <-time.After(RefreshInterval):
				refreshRoutes()
			}
		}
	}()

	go func() {
		for {
			select {
			case <-r.exit:
				return
			default:
				logger.Tracef("Router starting registry watch")
				w, err := r.options.Registry.Watch(registry.WatchDomain(registry.WildcardDomain))
				if err != nil {
					if logger.V(logger.DebugLevel, logger.DefaultLogger) {
						logger.Debugf("failed creating registry watcher: %v", err)
					}
					time.Sleep(time.Second)
					// in the event of an error reload routes
					refreshRoutes()
					continue
				}

				// watchRegistry calls stop when it's done
				if err := r.watchRegistry(w); err != nil {
					if logger.V(logger.DebugLevel, logger.DefaultLogger) {
						logger.Debugf("Error watching the registry: %v", err)
					}
					time.Sleep(time.Second)
					// in the event of an error reload routes
					refreshRoutes()
				}
			}
		}
	}()

	return nil
}

// Watch routes
func (r *registryRouter) Watch(opts ...router.WatchOption) (router.Watcher, error) {
	return r.table.Watch(opts...)
}

// String prints debugging information about router
func (r *registryRouter) String() string {
	return "registry"
}
