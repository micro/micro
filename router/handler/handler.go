package handler

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/network/router"
	pb "github.com/micro/go-micro/network/router/proto"
)

type Router struct {
	Router router.Router
}

func (r *Router) Lookup(ctx context.Context, req *pb.LookupRequest, resp *pb.LookupResponse) error {
	query := router.NewQuery(
		router.QueryDestination(req.Query.Destination),
	)

	routes, err := r.Router.Table().Lookup(query)
	if err != nil {
		return fmt.Errorf("failed to lookup routes: %s", err)
	}

	var respRoutes []*pb.Route
	for _, route := range routes {
		respRoute := &pb.Route{
			Destination: route.Destination,
			Gateway:     route.Gateway,
			Router:      route.Router,
			Network:     route.Network,
			Metric:      int64(route.Metric),
		}
		respRoutes = append(respRoutes, respRoute)
	}

	resp.Routes = respRoutes

	return nil
}
