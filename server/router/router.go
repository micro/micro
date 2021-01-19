package router

import (
	"context"

	pb "github.com/micro/micro/v3/proto/router"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/router"
)

// Router implements router handler
type Router struct {
	Router router.Router
}

// Lookup looks up routes in the routing table and returns them
func (r *Router) Lookup(ctx context.Context, req *pb.LookupRequest, resp *pb.LookupResponse) error {
	// build the options
	var options []router.LookupOption
	if v := req.Options; v != nil {
		if len(v.Address) > 0 {
			options = append(options, router.LookupAddress(v.Address))
		}
		if len(v.Gateway) > 0 {
			options = append(options, router.LookupGateway(v.Gateway))
		}
		if len(v.Router) > 0 {
			options = append(options, router.LookupRouter(v.Router))
		}
		if len(v.Network) > 0 {
			options = append(options, router.LookupNetwork(v.Network))
		}
		if len(v.Link) > 0 {
			options = append(options, router.LookupLink(v.Link))
		}
	}

	routes, err := r.Router.Lookup(req.Service, options...)
	if err == router.ErrRouteNotFound {
		return errors.NotFound("router.Router.Lookup", err.Error())
	} else if err != nil {
		return errors.InternalServerError("router.Router.Lookup", "failed to lookup routes: %v", err)
	}

	respRoutes := make([]*pb.Route, 0, len(routes))
	for _, route := range routes {
		respRoute := &pb.Route{
			Service:  route.Service,
			Address:  route.Address,
			Gateway:  route.Gateway,
			Network:  route.Network,
			Router:   route.Router,
			Link:     route.Link,
			Metric:   route.Metric,
			Metadata: route.Metadata,
		}
		respRoutes = append(respRoutes, respRoute)
	}

	resp.Routes = respRoutes

	return nil
}

// Watch streams routing table events
func (r *Router) Watch(ctx context.Context, req *pb.WatchRequest, stream pb.Router_WatchStream) error {
	watcher, err := r.Router.Watch()
	if err != nil {
		return errors.InternalServerError("router.Router.Watch", "failed creating event watcher: %v", err)
	}
	defer watcher.Stop()
	defer stream.Close()

	for {
		event, err := watcher.Next()
		if err == router.ErrWatcherStopped {
			return errors.InternalServerError("router.Router.Watch", "watcher stopped")
		}

		if err != nil {
			return errors.InternalServerError("router.Router.Watch", "error watching events: %v", err)
		}

		route := &pb.Route{
			Service:  event.Route.Service,
			Address:  event.Route.Address,
			Gateway:  event.Route.Gateway,
			Network:  event.Route.Network,
			Router:   event.Route.Router,
			Link:     event.Route.Link,
			Metric:   event.Route.Metric,
			Metadata: event.Route.Metadata,
		}

		tableEvent := &pb.Event{
			Id:        event.Id,
			Type:      pb.EventType(event.Type),
			Timestamp: event.Timestamp.UnixNano(),
			Route:     route,
		}

		if err := stream.Send(tableEvent); err != nil {
			return err
		}
	}
}
