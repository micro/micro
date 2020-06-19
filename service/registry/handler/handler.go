package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/service"
	pb "github.com/micro/go-micro/v2/registry/service/proto"
)

type Registry struct {
	// service id
	ID string
	// the publisher
	Publisher micro.Publisher
	// internal registry
	Registry registry.Registry
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

	return r.Publisher.Publish(ctx, event)
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

	// authorize the request
	if err := authorizeDomainAccess(ctx, options.Domain); err != nil {
		return err
	}

	// get the services in the namespace
	services, err := r.Registry.GetService(req.Service, registry.GetDomain(options.Domain))
	if err == registry.ErrNotFound {
		return errors.NotFound("go.micro.registry", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// serialize the response
	rsp.Services = make([]*pb.Service, len(services))
	for i, srv := range services {
		rsp.Services[i] = service.ToProto(srv)
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
	opts = append(opts, registry.RegisterDomain(req.Options.Domain))

	// authorize the request
	if err := authorizeDomainAccess(ctx, domain); err != nil {
		return err
	}

	// register the service
	if err := r.Registry.Register(service.ToService(req), opts...); err != nil {
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
		domain = registry.DefaultDomain
	}

	// authorize the request
	if err := authorizeDomainAccess(ctx, domain); err != nil {
		return err
	}

	// deregister the service
	if err := r.Registry.Deregister(service.ToService(req), registry.DeregisterDomain(domain)); err != nil {
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
		domain = registry.DefaultDomain
	}

	// authorize the request
	if err := authorizeDomainAccess(ctx, domain); err != nil {
		return err
	}

	// list the services from the registry
	services, err := r.Registry.ListServices(registry.ListDomain(domain))
	if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// serialize the response
	rsp.Services = make([]*pb.Service, len(services))
	for i, srv := range services {
		rsp.Services[i] = service.ToProto(srv)
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
	if err := authorizeDomainAccess(ctx, domain); err != nil {
		return err
	}

	// setup the watcher
	watcher, err := r.Registry.Watch(registry.WatchService(req.Service), registry.WatchDomain(domain))
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
			Service: service.ToProto(next.Service),
		})
		if err != nil {
			return errors.InternalServerError("go.micro.registry", err.Error())
		}
	}
}

// authorizeDomainAccess will return a go-micro error if the context cannot access the given domain
func authorizeDomainAccess(ctx context.Context, domain string) error {
	acc, ok := auth.AccountFromContext(ctx)

	// accounts are always required so we can identify the caller. If auth is not configured, the noop
	// auth implementation will return a blank account with the default domain set, allowing the caller
	// access to all resources
	if !ok {
		return errors.Unauthorized("go.micro.registry", "An account is required")
	}

	// anyone can access the default domain
	if domain == registry.DefaultDomain {
		return nil
	}

	// the server can access all domains
	if acc.Issuer == registry.DefaultDomain {
		return nil
	}

	// ensure the account is requesing access to it's own domain
	if acc.Issuer != domain {
		return errors.Forbidden("go.micro.registry", "An account issued by %v is required", domain)
	}

	return nil
}
