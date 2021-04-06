package registry

import (
	"context"
	"time"

	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/registry"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/registry/util"
)

type Registry struct {
	// service id
	ID string
	// the event
	Event *service.Event
}

func ActionToEventType(action string) registry.EventType {
	switch action {
	case "create":
		return registry.Create
	case "delete":
		return registry.Delete
	default:
		return registry.Update
	}
}

func (r *Registry) publishEvent(action string, service *pb.Service) error {
	// TODO: timestamp should be read from received event
	// Right now registry.Result does not contain timestamp
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
	var options registry.GetOptions
	if req.Options != nil && len(req.Options.Domain) > 0 {
		options.Domain = req.Options.Domain
	} else {
		options.Domain = registry.DefaultDomain
	}

	// authorize the request. Non admins can also do this
	publicNS := namespace.Public(registry.DefaultDomain)
	if err := namespace.Authorize(ctx, options.Domain, "registry.Registry.GetService", publicNS); err != nil {
		return err
	}

	// get the services in the namespace
	services, err := registry.DefaultRegistry.GetService(req.Service, registry.GetDomain(options.Domain))
	if err == registry.ErrNotFound || len(services) == 0 {
		return errors.NotFound("registry.Registry.GetService", registry.ErrNotFound.Error())
	} else if err != nil {
		return errors.InternalServerError("registry.Registry.GetService", err.Error())
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
	var opts []registry.RegisterOption
	var domain string

	// parse the options
	if req.Options != nil && req.Options.Ttl > 0 {
		ttl := time.Duration(req.Options.Ttl) * time.Second
		opts = append(opts, registry.RegisterTTL(ttl))
	}
	if req.Options != nil && len(req.Options.Domain) > 0 {
		domain = req.Options.Domain
	} else {
		domain = registry.DefaultDomain
	}
	opts = append(opts, registry.RegisterDomain(domain))

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, domain, "registry.Registry.Register"); err != nil {
		return err
	}

	// register the service
	if err := registry.DefaultRegistry.Register(util.ToService(req), opts...); err != nil {
		return errors.InternalServerError("registry.Registry.Register", err.Error())
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
		domain = registry.DefaultDomain
	}

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, domain, "registry.Registry.Deregister"); err != nil {
		return err
	}

	// deregister the service
	if err := registry.DefaultRegistry.Deregister(util.ToService(req), registry.DeregisterDomain(domain)); err != nil {
		return errors.InternalServerError("registry.Registry.Deregister", err.Error())
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
		domain = registry.DefaultDomain
	}

	// authorize the request. Non admins can also do this
	publicNS := namespace.Public(registry.DefaultDomain)
	if err := namespace.Authorize(ctx, domain, "registry.Registry.ListServices", publicNS); err != nil {
		return err
	}

	// list the services from the registry
	services, err := registry.DefaultRegistry.ListServices(registry.ListDomain(domain))
	if err != nil {
		return errors.InternalServerError("registry.Registry.ListServices", err.Error())
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
		domain = registry.DefaultDomain
	}

	// authorize the request
	publicNS := namespace.Public(registry.DefaultDomain)
	if err := namespace.Authorize(ctx, domain, "registry.Registry.Watch", publicNS); err != nil {
		return err
	}

	// setup the watcher
	watcher, err := registry.DefaultRegistry.Watch(registry.WatchService(req.Service), registry.WatchDomain(domain))
	if err != nil {
		return errors.InternalServerError("registry.Registry.Watch", err.Error())
	}

	for {
		next, err := watcher.Next()
		if err != nil {
			return errors.InternalServerError("registry.Registry.Watch", err.Error())
		}

		err = rsp.Send(&pb.Result{
			Action:  next.Action,
			Service: util.ToProto(next.Service),
		})
		if err != nil {
			return errors.InternalServerError("registry.Registry.Watch", err.Error())
		}
	}
}
