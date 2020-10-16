package client

import (
	"net/http"
	"sync"

	pb "github.com/micro/micro/v3/proto/router"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/router"
)

var (
	// name of the router service
	name = "router"
)

type svc struct {
	sync.RWMutex
	opts     router.Options
	callOpts []client.CallOption
	router   pb.RouterService
	table    *table
	exit     chan bool
	errChan  chan error
}

// NewRouter creates new service router and returns it
func NewRouter(opts ...router.Option) router.Router {
	// get default options
	options := router.DefaultOptions()

	// apply requested options
	for _, o := range opts {
		o(&options)
	}

	s := &svc{
		opts:   options,
		router: pb.NewRouterService(name, client.DefaultClient),
	}

	// set the router address to call
	if len(options.Address) > 0 {
		s.callOpts = []client.CallOption{
			client.WithAddress(options.Address),
			client.WithAuthToken(),
		}
	}
	// set the table
	s.table = &table{
		pb.NewTableService(name, client.DefaultClient),
		s.callOpts,
	}

	return s
}

// Init initializes router with given options
func (s *svc) Init(opts ...router.Option) error {
	s.Lock()
	defer s.Unlock()

	for _, o := range opts {
		o(&s.opts)
	}

	return nil
}

// Options returns router options
func (s *svc) Options() router.Options {
	s.Lock()
	opts := s.opts
	s.Unlock()

	return opts
}

// Table returns routing table
func (s *svc) Table() router.Table {
	return s.table
}

// Remote router cannot be closed
func (s *svc) Close() error {
	s.Lock()
	defer s.Unlock()

	select {
	case <-s.exit:
		return nil
	default:
		close(s.exit)
	}

	return nil
}

// Lookup looks up routes in the routing table and returns them
func (s *svc) Lookup(service string, q ...router.LookupOption) ([]router.Route, error) {
	// call the router
	query := router.NewLookup(q...)

	resp, err := s.router.Lookup(context.DefaultContext, &pb.LookupRequest{
		Service: service,
		Options: &pb.LookupOptions{
			Address: query.Address,
			Gateway: query.Gateway,
			Network: query.Network,
			Router:  query.Router,
			Link:    query.Link,
		},
	}, s.callOpts...)

	if verr := errors.FromError(err); verr != nil && verr.Code == http.StatusNotFound {
		return nil, router.ErrRouteNotFound
	} else if err != nil {
		return nil, err
	}

	routes := make([]router.Route, len(resp.Routes))
	for i, route := range resp.Routes {
		routes[i] = router.Route{
			Service:  route.Service,
			Address:  route.Address,
			Gateway:  route.Gateway,
			Network:  route.Network,
			Link:     route.Link,
			Metric:   route.Metric,
			Metadata: route.Metadata,
		}
	}

	return routes, nil
}

// Watch returns a watcher which allows to track updates to the routing table
func (s *svc) Watch(opts ...router.WatchOption) (router.Watcher, error) {
	rsp, err := s.router.Watch(context.DefaultContext, &pb.WatchRequest{}, s.callOpts...)
	if err != nil {
		return nil, err
	}
	options := router.WatchOptions{
		Service: "*",
	}
	for _, o := range opts {
		o(&options)
	}
	return newWatcher(rsp, options)
}

// Returns the router implementation
func (s *svc) String() string {
	return "service"
}
