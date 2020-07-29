package server

import (
	"context"
	"time"

	"github.com/micro/go-micro/v3/errors"
	log "github.com/micro/go-micro/v3/logger"
	goregistry "github.com/micro/go-micro/v3/registry"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/registry"
	pb "github.com/micro/micro/v3/service/registry/proto"
	"github.com/micro/micro/v3/service/registry/util"
)

type Registry struct {
	// service id
	ID string
	// the event
	Event *service.Event
}

func ActionToEventType(action string) goregistry.EventType {
	switch action {
	case "create":
		return goregistry.Create
	case "delete":
		return goregistry.Delete
	default:
		return goregistry.Update
	}
}

func (r *Registry) publishEvent(action string, service *pb.Service) error {
	// TODO: timestamp should be read from received event
	// Right now goregistry.Result does not contain timestamp
	event := &pb.Event{
		Id:        r.ID,
		Type:      pb.EventType(ActionToEventType(action)),
		Timestamp: time.Now().UnixNano(),
		Service:   service,
	}

	log.Debugf("publishing event %s for action %s", event.Id, action)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Event.Publish(ctx, event)
}

// GetService from the registry with the name requested
func (r *Registry) GetService(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	// parse the options
	var options goregistry.GetOptions
	if req.Options != nil && len(req.Options.Domain) > 0 {
		options.Domain = req.Options.Domain
	} else {
		options.Domain = goregistry.DefaultDomain
	}

	// authorize the request
	publicNS := namespace.Public(goregistry.DefaultDomain)
	if err := namespace.Authorize(ctx, options.Domain, publicNS); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.registry", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.registry", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// get the services in the namespace
	services, err := registry.GetService(req.Service, goregistry.GetDomain(options.Domain))
	if err == goregistry.ErrNotFound {
		return errors.NotFound("go.micro.registry", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// serialize the response
	rsp.Services = make([]*pb.Service, len(services))
	for i, srv := range services {
		rsp.Services[i] = util.ToProto(srv)
	}

	return nil
}

// Register a service
func (r *Registry) Register(ctx context.Context, req *pb.Service, rsp *pb.EmptyResponse) error {
	var opts []goregistry.RegisterOption
	var domain string

	// parse the options
	if req.Options != nil && req.Options.Ttl > 0 {
		ttl := time.Duration(req.Options.Ttl) * time.Second
		opts = append(opts, goregistry.RegisterTTL(ttl))
	}
	if req.Options != nil && len(req.Options.Domain) > 0 {
		domain = req.Options.Domain
	} else {
		domain = goregistry.DefaultDomain
	}
	opts = append(opts, goregistry.RegisterDomain(domain))

	// authorize the request
	if err := namespace.Authorize(ctx, domain); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.registry", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.registry", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// register the service
	if err := registry.Register(util.ToService(req), opts...); err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// publish the event
	go r.publishEvent("create", req)

	return nil
}

// Deregister a service
func (r *Registry) Deregister(ctx context.Context, req *pb.Service, rsp *pb.EmptyResponse) error {
	// parse the options
	var domain string
	if req.Options != nil && len(req.Options.Domain) > 0 {
		domain = req.Options.Domain
	} else {
		domain = goregistry.DefaultDomain
	}

	// authorize the request
	if err := namespace.Authorize(ctx, domain); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.registry", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.registry", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// deregister the service
	if err := registry.Deregister(util.ToService(req), goregistry.DeregisterDomain(domain)); err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// publish the event
	go r.publishEvent("delete", req)

	return nil
}

// ListServices returns all the services
func (r *Registry) ListServices(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// parse the options
	var domain string
	if req.Options != nil && len(req.Options.Domain) > 0 {
		domain = req.Options.Domain
	} else {
		domain = goregistry.DefaultDomain
	}

	// authorize the request
	publicNS := namespace.Public(goregistry.DefaultDomain)
	if err := namespace.Authorize(ctx, domain, publicNS); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.registry", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.registry", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// list the services from the registry
	services, err := registry.ListServices(goregistry.ListDomain(domain))
	if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// serialize the response
	rsp.Services = make([]*pb.Service, len(services))
	for i, srv := range services {
		rsp.Services[i] = util.ToProto(srv)
	}

	return nil
}

// Watch a service for changes
func (r *Registry) Watch(ctx context.Context, req *pb.WatchRequest, rsp pb.Registry_WatchStream) error {
	// parse the options
	var domain string
	if req.Options != nil && len(req.Options.Domain) > 0 {
		domain = req.Options.Domain
	} else {
		domain = goregistry.DefaultDomain
	}

	// authorize the request
	publicNS := namespace.Public(goregistry.DefaultDomain)
	if err := namespace.Authorize(ctx, domain, publicNS); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.registry", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.registry", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// setup the watcher
	watcher, err := registry.Watch(goregistry.WatchService(req.Service), goregistry.WatchDomain(domain))
	if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	for {
		next, err := watcher.Next()
		if err != nil {
			return errors.InternalServerError("go.micro.registry", err.Error())
		}

		err = rsp.Send(&pb.Result{
			Action:  next.Action,
			Service: util.ToProto(next.Service),
		})
		if err != nil {
			return errors.InternalServerError("go.micro.registry", err.Error())
		}
	}
}
