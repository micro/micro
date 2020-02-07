package main

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/util/log"

	pb "github.com/micro/go-micro/v2/platform/api/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.platform"),
		micro.Version("latest"),
	)
	service.Init()

	pb.RegisterPlatformHandler(service.Server(), NewHandler())

	if err := service.Run(); err != nil {
		log.Error(err)
	}
}

// Handler is an impementation of the platform api
type Handler struct {
	runtime runtime.Runtime
}

// NewHandler returns an initialized Handler
func NewHandler() *Handler {
	return &Handler{runtime: runtime.DefaultRuntime}
}

// CreateService deploys a service on the platform
func (h *Handler) CreateService(ctx context.Context, req *pb.CreateServiceRequest, rsp *pb.CreateServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	return h.runtime.Create(deserializeService(req.Service))
}

// ReadService returns information about services matching the query
func (h *Handler) ReadService(ctx context.Context, req *pb.ReadServiceRequest, rsp *pb.ReadServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	resp, err := h.runtime.Read(
		runtime.ReadType(req.Service.Type),
		runtime.ReadService(req.Service.Name),
		runtime.ReadVersion(req.Service.Version),
	)
	if err != nil {
		return err
	}

	rsp.Services = make([]*pb.Service, len(resp))
	for i, s := range resp {
		rsp.Services[i] = serializeService(s)
	}

	return nil
}

// UpdateService updates a service running on the platform
func (h *Handler) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest, rsp *pb.UpdateServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	return h.runtime.Update(deserializeService(req.Service))
}

// DeleteService terminates a service running on the platform
func (h *Handler) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest, rsp *pb.DeleteServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	return h.runtime.Delete(deserializeService(req.Service))
}

// ListServices returns all the services running on the platform
func (h *Handler) ListServices(ctx context.Context, req *pb.ListServicesRequest, rsp *pb.ListServicesResponse) error {
	resp, err := h.runtime.List()
	if err != nil {
		return err
	}

	rsp.Services = make([]*pb.Service, len(resp))
	for i, s := range resp {
		rsp.Services[i] = serializeService(s)
	}

	return nil
}

func serializeService(srv *runtime.Service) *pb.Service {
	return &pb.Service{
		Name:     srv.Name,
		Version:  srv.Version,
		Source:   srv.Source,
		Metadata: srv.Metadata,
	}
}

func deserializeService(srv *pb.Service) *runtime.Service {
	return &runtime.Service{
		Name:     srv.Name,
		Version:  srv.Version,
		Source:   srv.Source,
		Metadata: srv.Metadata,
	}
}
