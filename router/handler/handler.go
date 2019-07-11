package handler

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/network/router"
	pb "github.com/micro/go-micro/network/router/proto"
	"github.com/micro/go-micro/network/router/table"
)

// Router implements router handler
type Router struct {
	Router router.Router
}

// Lookup looks up routes in the routing table and returns them
func (r *Router) Lookup(ctx context.Context, req *pb.LookupRequest, resp *pb.LookupResponse) error {
	query := table.NewQuery(
		table.QueryService(req.Query.Service),
	)

	routes, err := r.Router.Lookup(query)
	if err != nil {
		return fmt.Errorf("failed to lookup routes: %s", err)
	}

	var respRoutes []*pb.Route
	for _, route := range routes {
		respRoute := &pb.Route{
			Service: route.Service,
			Address: route.Address,
			Gateway: route.Gateway,
			Network: route.Network,
			Link:    route.Link,
			Metric:  int64(route.Metric),
		}
		respRoutes = append(respRoutes, respRoute)
	}

	resp.Routes = respRoutes

	return nil
}

// List returns all routes in the routing table
func (r *Router) List(ctx context.Context, req *pb.ListRequest, resp *pb.ListResponse) error {
	routes, err := r.Router.List()
	if err != nil {
		return fmt.Errorf("failed to list routes: %s", err)
	}

	var respRoutes []*pb.Route
	for _, route := range routes {
		respRoute := &pb.Route{
			Service: route.Service,
			Address: route.Address,
			Gateway: route.Gateway,
			Network: route.Network,
			Link:    route.Link,
			Metric:  int64(route.Metric),
		}
		respRoutes = append(respRoutes, respRoute)
	}

	resp.Routes = respRoutes

	return nil
}

// Watch streans routing table events
func (r *Router) Watch(ctx context.Context, req *pb.WatchRequest, stream pb.Router_WatchStream) error {
	watcher, err := r.Router.Watch()
	if err != nil {
		return fmt.Errorf("failed creating event watcher: %v", err)
	}

	defer stream.Close()

	for {
		event, err := watcher.Next()
		if err == table.ErrWatcherStopped {
			break
		}

		if err != nil {
			return fmt.Errorf("error watching events: %s", err)
		}

		route := &pb.Route{
			Service: event.Route.Service,
			Address: event.Route.Address,
			Gateway: event.Route.Gateway,
			Network: event.Route.Network,
			Link:    event.Route.Link,
			Metric:  int64(event.Route.Metric),
		}

		tableEvent := &pb.Event{
			Type:      pb.EventType(event.Type),
			Timestamp: event.Timestamp.UnixNano(),
			Route:     route,
		}

		if err := stream.Send(tableEvent); err != nil {
			return err
		}
	}

	return nil
}
