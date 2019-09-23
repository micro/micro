// Package registry is the micro registry
package registry

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/handler"
	pb "github.com/micro/go-micro/registry/proto"
	"github.com/micro/go-micro/registry/service"
	"github.com/micro/go-micro/util/log"
	rcli "github.com/micro/micro/cli"
)

var (
	// Name of the registry
	Name = "go.micro.registry"
	// The address of the registry
	Address = ":8000"
	// Topic to publish registry events to
	Topic = "go.micro.registry.events"
)

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

// Sub processes registry events
type sub struct {
	// id is registry id
	id string
	// registry is service registry
	registry registry.Registry
}

// Process processes registry events
func (s *sub) Process(ctx context.Context, event *pb.Event) error {
	log.Logf("[registry] received %s event from: %s", registry.EventType(event.Type), event.Id)
	if event.Id == s.id {
		log.Logf("[registry] skipping event")
		return nil
	}

	// decode protobuf to registry.Service
	svc := service.ToService(event.Service)

	switch registry.EventType(event.Type) {
	case registry.Create, registry.Update:
		if err := s.registry.Register(svc); err != nil {
			log.Logf("[registry] failed to register service: %s", svc.Name)
			return err
		}
	case registry.Delete:
		if err := s.registry.Deregister(svc); err != nil {
			log.Logf("[registry] failed to deregister service: %s", svc.Name)
			return err
		}
	}

	return nil
}

// reg is micro registry
type reg struct {
	// id is registry id
	id string
	// registry is micro registry
	registry.Registry
	// publisher to publish registry events
	micro.Publisher
}

// newRegsitry creates new micro registry and returns it
func newRegistry(service micro.Service, registry registry.Registry) *reg {
	id := uuid.New().String()
	s := &sub{
		id:       id,
		registry: registry,
	}

	// register subscriber
	if err := micro.RegisterSubscriber(Topic, service.Server(), s); err != nil {
		log.Logf("[registry] failed to subscribe to events: %s", err)
		os.Exit(1)
	}

	return &reg{
		id:        id,
		Registry:  registry,
		Publisher: micro.NewPublisher(Topic, service.Client()),
	}
}

// Publish publishes registry events to other registries to consume
func (r *reg) PublishEvents(w registry.Watcher) error {
	defer w.Stop()

	var watchErr error

	for {
		res, err := w.Next()
		if err != nil {
			if err != registry.ErrWatcherStopped {
				watchErr = err
			}
			break
		}

		// encode *registry.Service into protobuf messag
		svc := service.ToProto(res.Service)

		// TODO: timestamp should be read from received event
		// Right now registry.Result does not contain timestamp
		event := &pb.Event{
			Id:        r.id,
			Type:      pb.EventType(ActionToEventType(res.Action)),
			Timestamp: time.Now().UnixNano(),
			Service:   svc,
		}

		if err := r.Publish(context.Background(), event); err != nil {
			log.Logf("[registry] error publishing event: %v", err)
			return fmt.Errorf("error publishing event: %v", err)
		}
	}

	return watchErr
}

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
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
	if i := time.Duration(ctx.GlobalInt("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.GlobalInt("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// set address
	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// new service
	service := micro.NewService(srvOpts...)

	pb.RegisterRegistryHandler(service.Server(), &handler.Registry{
		// using the mdns registry
		Registry: service.Options().Registry,
	})

	reg := newRegistry(service, service.Options().Registry)

	// create registry watcher
	watcher, err := service.Options().Registry.Watch()
	if err != nil {
		log.Logf("[registry] failed creating watcher: %v", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	// error channel to collect registry errors
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- reg.PublishEvents(watcher)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- service.Run()
	}()

	// we block here until either service or server fails
	if err := <-errChan; err != nil {
		log.Logf("[registry] error running the registry: %v", err)
		if err != registry.ErrWatcherStopped {
			watcher.Stop()
			os.Exit(1)
		}
		os.Exit(1)
	}

	// stop registry watcher
	watcher.Stop()

	wg.Wait()

	log.Logf("[registry] successfully stopped")
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "registry",
		Usage: "Run the service registry",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the registry http address e.g 0.0.0.0:8080",
				EnvVar: "MICRO_REGISTRY_ADDRESS",
			},
		},
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
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

	return []cli.Command{command}
}
