package server

import (
	"context"
	"time"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/micro/v2/service"
	mureg "github.com/micro/micro/v2/service/registry"
	pb "github.com/micro/micro/v2/service/registry/proto"
	"github.com/micro/micro/v2/service/registry/util"
)

var (
	// name of the registry
	name = "go.micro.registry"
	// address of the registry
	address = ":8000"
	// topic to publish registry events to
	topic = "go.micro.registry.events"
)

// Sub processes registry events
type subscriber struct {
	// id is registry id
	Id string
	// registry is service registry
	Registry registry.Registry
}

// Process processes registry events
func (s *subscriber) Process(ctx context.Context, event *pb.Event) error {
	if event.Id == s.Id {
		log.Tracef("skipping own %s event: %s for: %s", registry.EventType(event.Type), event.Id, event.Service.Name)
		return nil
	}

	log.Debugf("received %s event from: %s for: %s", registry.EventType(event.Type), event.Id, event.Service.Name)

	// no service
	if event.Service == nil {
		return nil
	}

	// decode protobuf to registry.Service
	svc := util.ToService(event.Service)

	// default ttl to 1 minute
	ttl := time.Minute

	// set ttl if it exists
	if opts := event.Service.Options; opts != nil {
		if opts.Ttl > 0 {
			ttl = time.Duration(opts.Ttl) * time.Second
		}
	}

	switch registry.EventType(event.Type) {
	case registry.Create, registry.Update:
		log.Debugf("registering service: %s", svc.Name)
		if err := s.Registry.Register(svc, registry.RegisterTTL(ttl)); err != nil {
			log.Debugf("failed to register service: %s", svc.Name)
			return err
		}
	case registry.Delete:
		log.Debugf("deregistering service: %s", svc.Name)
		if err := s.Registry.Deregister(svc); err != nil {
			log.Debugf("failed to deregister service: %s", svc.Name)
			return err
		}
	}

	return nil
}

func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}

	// service opts
	srvOpts := []service.Option{service.Name(name)}
	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterInterval(i*time.Second))
	}

	// set address
	if len(address) > 0 {
		srvOpts = append(srvOpts, service.Address(address))
	}

	// new service
	srv := service.New(srvOpts...)
	// get server id
	id := srv.Server().Options().Id

	// register the handler
	pb.RegisterRegistryHandler(srv.Server(), &Registry{
		ID:       id,
		Event:    service.NewEvent(topic),
		Registry: mureg.DefaultRegistry,
	})

	// run the service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
