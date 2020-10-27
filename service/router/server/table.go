package server

import (
	"context"

	pb "github.com/micro/micro/v3/proto/router"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/router"
)

type Table struct {
	Router router.Router
}

func (t *Table) Create(ctx context.Context, route *pb.Route, resp *pb.CreateResponse) error {
	err := t.Router.Table().Create(router.Route{
		Service:  route.Service,
		Address:  route.Address,
		Gateway:  route.Gateway,
		Network:  route.Network,
		Router:   route.Router,
		Link:     route.Link,
		Metric:   route.Metric,
		Metadata: route.Metadata,
	})
	if err != nil {
		return errors.InternalServerError("router.Table.Create", "failed to create route: %s", err)
	}

	return nil
}

func (t *Table) Update(ctx context.Context, route *pb.Route, resp *pb.UpdateResponse) error {
	err := t.Router.Table().Update(router.Route{
		Service:  route.Service,
		Address:  route.Address,
		Gateway:  route.Gateway,
		Network:  route.Network,
		Router:   route.Router,
		Link:     route.Link,
		Metric:   route.Metric,
		Metadata: route.Metadata,
	})
	if err != nil {
		return errors.InternalServerError("router.Table.Update", "failed to update route: %s", err)
	}

	return nil
}

func (t *Table) Delete(ctx context.Context, route *pb.Route, resp *pb.DeleteResponse) error {
	err := t.Router.Table().Delete(router.Route{
		Service:  route.Service,
		Address:  route.Address,
		Gateway:  route.Gateway,
		Network:  route.Network,
		Router:   route.Router,
		Link:     route.Link,
		Metric:   route.Metric,
		Metadata: route.Metadata,
	})
	if err != nil {
		return errors.InternalServerError("route.Table.Delete", "failed to delete route: %s", err)
	}

	return nil
}

func (t *Table) Read(ctx context.Context, req *pb.ReadRequest, resp *pb.ReadResponse) error {
	routes, err := t.Router.Table().Read(router.ReadService(req.Service))
	if err != nil {
		return errors.InternalServerError("router.Table.Read", "failed to lookup routes: %s", err)
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
