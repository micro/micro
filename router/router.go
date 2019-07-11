package router

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/network/router"
	pb "github.com/micro/go-micro/network/router/proto"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/router/handler"
)

var (
	// Name of the router microservice
	Name = "go.micro.router"
	// Address is the router microservice bind address
	Address = ":8084"
	// Router is the router gossip bind address
	Router = ":9093"
	// Network is the network id
	Network = "local"
	// Topic is router events topic
	Topic = "go.micro.router.events"
)

// Pub publishes router events
type Pub struct {
	micro.Publisher
}

// NewPub creates new publisher and returns it
func NewPub(topic string, client client.Client) *Pub {
	return &Pub{
		Publisher: micro.NewPublisher(Topic, client),
	}
}

// PubEvents publishes advertised events
func (p *Pub) PubEvents(ch <-chan *router.Advert) error {
	for advert := range ch {
		for _, event := range advert.Events {
			route := &pb.Route{
				Service: event.Route.Service,
				Address: event.Route.Address,
				Gateway: event.Route.Gateway,
				Network: event.Route.Network,
				Link:    event.Route.Link,
				Metric:  int64(event.Route.Metric),
			}
			event := &pb.TableEvent{
				Type:      pb.EventType(event.Type),
				Timestamp: event.Timestamp.UnixNano(),
				Route:     route,
			}

			if err := p.Publish(context.Background(), event); err != nil {
				log.Logf("[router] error publishing event: %v", err)
				return fmt.Errorf("error publishing event: %v", err)
			}
		}
	}

	return nil
}

// Sub processes router events
type Sub struct{}

// Process advertised events
func (s *Sub) Process(ctx context.Context, event *pb.TableEvent) error {
	// TODO: filter the events which were originated by you
	log.Logf("[router] Received event: %+v", event)
	return nil
}

// run runs the micro server
func run(ctx *cli.Context, srvOpts ...micro.Option) {
	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("router_address")) > 0 {
		Router = ctx.String("router")
	}
	if len(ctx.String("network_address")) > 0 {
		Network = ctx.String("network")
	}
	// default gateway address
	var gateway string
	if len(ctx.String("gateway_address")) > 0 {
		gateway = ctx.String("gateway")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	r := router.NewRouter(
		router.Id(service.Server().Options().Id),
		router.Address(Router),
		router.Network(Network),
		router.Registry(service.Client().Options().Registry),
		router.Gateway(gateway),
	)

	// register router handler
	pb.RegisterRouterHandler(
		service.Server(),
		&handler.Router{Router: r},
	)

	log.Log("[router] starting to advertise")

	advertChan, err := r.Advertise()
	if err != nil {
		log.Logf("[router] failed to start: %s", err)
		os.Exit(1)
	}

	// create event publisher
	pub := NewPub(Topic, service.Client())

	// register subscriber
	if err := micro.RegisterSubscriber(Topic, service.Server(), new(Sub)); err != nil {
		log.Logf("[router] failed to register subscriber: %s", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- pub.PubEvents(advertChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- service.Run()
	}()

	// we block here until either service or server fails
	if err := <-errChan; err != nil {
		log.Logf("[router] error running the router: %v", err)
	}

	log.Log("[router] attempting to stop the router")

	// stop the router
	if err := r.Stop(); err != nil {
		log.Logf("[router] failed to stop: %s", err)
		os.Exit(1)
	}

	wg.Wait()

	log.Logf("[router] successfully stopped")
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "router",
		Usage: "Run the micro network router",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "router_address",
				Usage:  "Set the micro router address :9093",
				EnvVar: "MICRO_ROUTER_ADDRESS",
			},
			cli.StringFlag{
				Name:   "network_address",
				Usage:  "Set the micro network address: local",
				EnvVar: "MICRO_NETWORK_ADDRESS",
			},
			cli.StringFlag{
				Name:   "gateway_address",
				Usage:  "Set the micro default gateway address :9094",
				EnvVar: "MICRO_GATEWAY_ADDRESS",
			},
		},
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
		},
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
