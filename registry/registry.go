// Package registry is the micro registry
package registry

import (
	"context"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/service"
	pb "github.com/micro/go-micro/v2/registry/service/proto"
	rcli "github.com/micro/micro/v2/cli"
	"github.com/micro/micro/v2/registry/handler"
)

var (
	// Name of the registry
	Name = "go.micro.registry"
	// The address of the registry
	Address = ":8000"
	// Topic to publish registry events to
	Topic = "go.micro.registry.events"
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
	svc := service.ToService(event.Service)

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

func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Info("registry")

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// service opts
	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// set address
	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// new service
	service := micro.NewService(srvOpts...)
	// get server id
	id := service.Server().Options().Id

	/*
		// create the subscriber
		s := &subscriber{
			Id:       id,
			Registry: service.Options().Registry,
		}

		// register the subscriber
			if err := micro.RegisterSubscriber(Topic, service.Server(), s); err != nil {
				log.Debugf("failed to subscribe to events: %s", err)
				os.Exit(1)
			}
	*/

	// register the handler
	pb.RegisterRegistryHandler(service.Server(), &handler.Registry{
		Id:        id,
		Publisher: micro.NewPublisher(Topic, service.Client()),
		Registry:  service.Options().Registry,
	})

	// run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "registry",
		Usage: "Run the service registry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the registry http address e.g 0.0.0.0:8000",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			},
		},
		Action: func(ctx *cli.Context) error {
			Run(ctx, options...)
			return nil
		},
		Subcommands: rcli.RegistryCommands(),
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []*cli.Command{command}
}
