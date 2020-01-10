package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/runtime"
	pb "github.com/micro/go-micro/runtime/service/proto"
)

type Runtime struct {
	// The runtime used to manage services
	Runtime runtime.Runtime
	// The client used to publish events
	Client micro.Publisher
}

func (r *Runtime) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	var options []runtime.CreateOption
	if req.Options != nil {
		options = toCreateOptions(req.Options)
	}

	service := toService(req.Service)
	err := r.Runtime.Create(service, options...)
	if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the delete event
	r.Client.Publish(ctx, &pb.Event{
		Type:      "create",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

func (r *Runtime) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	var options []runtime.ReadOption

	if req.Options != nil {
		options = toReadOptions(req.Options)
	}

	services, err := r.Runtime.Read(options...)
	if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	for _, service := range services {
		rsp.Services = append(rsp.Services, toProto(service))
	}

	return nil
}

func (r *Runtime) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// TODO: add opts
	service := toService(req.Service)
	err := r.Runtime.Update(service)
	if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the delete event
	r.Client.Publish(ctx, &pb.Event{
		Type:      "update",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

func (r *Runtime) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// TODO: add opts
	service := toService(req.Service)
	err := r.Runtime.Delete(service)
	if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the delete event
	r.Client.Publish(ctx, &pb.Event{
		Type:      "delete",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

func (r *Runtime) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	services, err := r.Runtime.List()
	if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	for _, service := range services {
		rsp.Services = append(rsp.Services, toProto(service))
	}

	return nil
}
