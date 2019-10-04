// Package registry is the micro registry
package registry

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
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
	// SyncTime defines time interval to periodically sync registries
	SyncTime = 5 * time.Second
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
	log.Debugf("[registry] received %s event from: %s for: %s", registry.EventType(event.Type), event.Id, event.Service.Name)
	if event.Id == s.id {
		log.Debugf("[registry] skipping event")
		return nil
	}

	// decode protobuf to registry.Service
	svc := service.ToService(event.Service)

	switch registry.EventType(event.Type) {
	case registry.Create, registry.Update:
		log.Debugf("[registry] registering service: %s", svc.Name)
		if err := s.registry.Register(svc); err != nil {
			log.Debugf("[registry] failed to register service: %s", svc.Name)
			return err
		}
	case registry.Delete:
		log.Debugf("[registry] deregistering service: %s", svc.Name)
		if err := s.registry.Deregister(svc); err != nil {
			log.Debugf("[registry] failed to deregister service: %s", svc.Name)
			return err
		}
	}

	return nil
}

// reg is micro registry
type reg struct {
	// registry is micro registry
	registry.Registry
	// id is registry id
	id string
	// client is service client
	client client.Client
	// exit stops the registry
	exit chan bool
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
		log.Debugf("[registry] failed to subscribe to events: %s", err)
		os.Exit(1)
	}

	return &reg{
		Registry: registry,
		id:       id,
		client:   service.Client(),
		exit:     make(chan bool),
	}
}

// Publish publishes registry events to other registries to consume
func (r *reg) PublishEvents(w registry.Watcher) error {
	defer w.Stop()

	// create a publisher
	p := micro.NewPublisher(Topic, r.client)
	// track watcher errors
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

		log.Debugf("[registry] publishing event %s for action %s", event.Id, res.Action)

		select {
		case <-r.exit:
			return nil
		default:
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := p.Publish(ctx, event); err != nil {
				log.Debugf("[registry] error publishing event: %v", err)
				return fmt.Errorf("error publishing event: %v", err)
			}
		}
	}

	return watchErr
}

func (r *reg) syncRecords(nodes []string) error {
	if len(nodes) == 0 {
		log.Debugf("[registry] no nodes to sync with. skipping")
		return nil
	}

	c := pb.NewRegistryService(Name, r.client)
	resp, err := c.ListServices(context.Background(), &pb.ListRequest{}, client.WithAddress(nodes...))
	if err != nil {
		log.Debugf("[registry] failed sync: %v", err)
		return err
	}

	for _, pbService := range resp.Services {
		svc := service.ToService(pbService)
		log.Debugf("[registry] registering service: %s", svc.Name)
		if err := r.Register(svc); err != nil {
			log.Debugf("[registry] failed to register service: %v", svc.Name)
			return err
		}
	}

	return nil
}

func (r *reg) Sync(nodes []string) error {
	sync := time.NewTicker(SyncTime)
	defer sync.Stop()

	for {
		select {
		case <-r.exit:
			return nil
		case <-sync.C:
			log.Debugf("[registry] syncing records")
			if err := r.syncRecords(nodes); err != nil {
				log.Debugf("[registry] failed to sync registry records: %v", err)
				continue
			}
		}
	}

	return nil
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
	var nodes []string
	if len(ctx.String("nodes")) > 0 {
		nodes = strings.Split(ctx.String("nodes"), ",")
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
		log.Debugf("[registry] failed creating watcher: %v", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- reg.PublishEvents(watcher)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- reg.Sync(nodes)
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
			close(reg.exit)
			os.Exit(1)
		}
		os.Exit(1)
	}

	// stop registry watcher
	watcher.Stop()
	close(reg.exit)

	wg.Wait()

	log.Debugf("[registry] successfully stopped")
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "registry",
		Usage: "Run the service registry",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the registry http address e.g 0.0.0.0:8000",
				EnvVar: "MICRO_SERVER_ADDRESS",
			},
			cli.StringFlag{
				Name:   "nodes",
				Usage:  "Set the micro registry nodes to connect to. This can be a comma separated list.",
				EnvVar: "MICRO_REGISTRY_NODES",
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
