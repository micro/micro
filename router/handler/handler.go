package handler

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/network/router"
	pb "github.com/micro/go-micro/network/router/proto"
	"github.com/micro/go-micro/network/router/table"
)

type Router struct {
	Router router.Router
}

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
