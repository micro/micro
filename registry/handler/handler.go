package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/service"
	pb "github.com/micro/go-micro/v2/registry/service/proto"
	"github.com/micro/micro/v2/internal/namespace"
)

type Registry struct {
	// service id
	Id string
	// the publisher
	Publisher micro.Publisher
	// internal registry
	Registry registry.Registry
	// auth to verify clients
	Auth auth.Auth
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
		Id:        r.Id,
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
	// get the services in the requested namespace, e.g. the "foo" namespace. name
	// includes the namespace as the prefix, e.g. 'foo/go.micro.service.bar'
	name := namespace.FromContext(ctx) + service.NameSeperator + req.Service
	services, err := r.Registry.GetService(name)
	if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// get the services in the default namespace if this wasn't the namespace
	// requested.
	if namespace.FromContext(ctx) != namespace.DefaultNamespace {
		name := namespace.DefaultNamespace + service.NameSeperator + req.Service
		defaultServices, err := r.Registry.GetService(name)
		if err != nil {
			return errors.InternalServerError("go.micro.registry", err.Error())
		}
		services = append(services, defaultServices...)
	}

	// serialize the services. service.ToProto will remove the namespace from
	// the service name so 'foo/go.micro.service.bar' will become just 'go.micro.service.bar'.
	for _, srv := range services {
		rsp.Services = append(rsp.Services, service.ToProto(srv))
	}
	return nil
}

// Register a service
func (r *Registry) Register(ctx context.Context, req *pb.Service, rsp *pb.EmptyResponse) error {
	var regOpts []registry.RegisterOption
	if req.Options != nil {
		ttl := time.Duration(req.Options.Ttl) * time.Second
		regOpts = append(regOpts, registry.RegisterTTL(ttl))
	}

	service := service.ToService(req, service.WithNamespace(namespace.FromContext(ctx)))
	if err := r.Registry.Register(service, regOpts...); err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// publish the event
	go r.publishEvent("create", req)

	return nil
}

// Deregister a service
func (r *Registry) Deregister(ctx context.Context, req *pb.Service, rsp *pb.EmptyResponse) error {
	service := service.ToService(req, service.WithNamespace(namespace.FromContext(ctx)))
	if err := r.Registry.Deregister(service); err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	// publish the event
	go r.publishEvent("delete", req)

	return nil
}

// ListServices returns all the services
func (r *Registry) ListServices(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	fmt.Println(namespace.FromContext(ctx))

	services, err := r.Registry.ListServices()
	if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	for _, srv := range services {
		// check to see if the service belongs to the defaut namespace
		// or the contexts namespace. TODO: think about adding a prefix
		//argument to ListServices
		if !canReadService(ctx, srv) {
			continue
		}

		rsp.Services = append(rsp.Services, service.ToProto(srv))
	}

	return nil
}

// Watch a service for changes
func (r *Registry) Watch(ctx context.Context, req *pb.WatchRequest, rsp pb.Registry_WatchStream) error {
	watcher, err := r.Registry.Watch(registry.WatchService(req.Service))
	if err != nil {
		return errors.InternalServerError("go.micro.registry", err.Error())
	}

	for {
		next, err := watcher.Next()
		if err != nil {
			return errors.InternalServerError("go.micro.registry", err.Error())
		}
		if !canReadService(ctx, next.Service) {
			continue
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

// canReadService is a helper function which returns a boolean indicating
// if a context can read a service.
func canReadService(ctx context.Context, srv *registry.Service) bool {
	// all users can read from the default namespace
	if strings.HasPrefix(srv.Name, namespace.DefaultNamespace+service.NameSeperator) {
		return true
	}

	// the service belongs to the contexts namespace
	if strings.HasPrefix(srv.Name, namespace.FromContext(ctx)+service.NameSeperator) {
		return true
	}

	return false
}
