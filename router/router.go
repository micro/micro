package router

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/network/router"
	pb "github.com/micro/go-micro/network/router/proto"
	"github.com/micro/go-micro/network/router/table"
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
	// Topic is router adverts topic
	Topic = "go.micro.router.adverts"
)

// Sub processes router events
type sub struct {
	router router.Router
}

// Process processes router adverts
func (s *sub) Process(ctx context.Context, advert *pb.Advert) error {
	log.Logf("[router] received advert: %+v", advert)
	if advert.Id == s.router.Options().Id {
		log.Logf("[router] skipping advert")
		return nil
	}

	for _, event := range advert.Events {
		action := fmt.Sprintf("%s", table.EventType(event.Type))
		route := table.Route{
			Service: event.Route.Service,
			Address: event.Route.Address,
			Gateway: event.Route.Gateway,
			Network: event.Route.Network,
			Link:    event.Route.Link,
			Metric:  int(event.Route.Metric),
		}

		switch action {
		case "create":
			if err := s.router.Create(route); err != nil && err != table.ErrDuplicateRoute {
				return fmt.Errorf("failed adding route for service %s: %s", route.Service, err)
			}
		case "update":
			if err := s.router.Update(route); err != nil && err != table.ErrDuplicateRoute {
				return fmt.Errorf("failed updating route for service %s: %s", route.Service, err)
			}
		case "delete":
			if err := s.router.Delete(route); err != nil && err != table.ErrRouteNotFound {
				return fmt.Errorf("failed deleting route for service %s: %s", route.Service, err)
			}
		default:
			return fmt.Errorf("failed to manage route for service %s. Unknown action: %s", route.Service, action)
		}
	}
	return nil
}

// rtr is micro router
type rtr struct {
	// router is the micro router
	router.Router
	// publisher to publish router adverts
	micro.Publisher
}

// newRouter creates new micro router and returns it
func newRouter(service micro.Service, router router.Router) *rtr {
	s := &sub{
		router: router,
	}

	// register subscriber
	if err := micro.RegisterSubscriber(Topic, service.Server(), s); err != nil {
		log.Logf("[router] failed to subscribe to adverts: %s", err)
		os.Exit(1)
	}

	return &rtr{
		Router:    router,
		Publisher: micro.NewPublisher(Topic, service.Client()),
	}
}

// PublishAdverts publishes adverts for other routers to consume
func (r *rtr) PublishAdverts(ch <-chan *router.Advert) error {
	for advert := range ch {
		var events []*pb.Event
		for _, event := range advert.Events {
			route := &pb.Route{
				Service: event.Route.Service,
				Address: event.Route.Address,
				Gateway: event.Route.Gateway,
				Network: event.Route.Network,
				Link:    event.Route.Link,
				Metric:  int64(event.Route.Metric),
			}
			e := &pb.Event{
				Type:      pb.EventType(event.Type),
				Timestamp: event.Timestamp.UnixNano(),
				Route:     route,
			}
			events = append(events, e)
		}

		a := &pb.Advert{
			Id:        r.Options().Id,
			Type:      pb.AdvertType(advert.Type),
			Timestamp: advert.Timestamp.UnixNano(),
			Events:    events,
		}

		if err := r.Publish(context.Background(), a); err != nil {
			log.Logf("[router] error publishing advert: %v", err)
			return fmt.Errorf("error publishing advert: %v", err)
		}
	}

	return nil
}

// stop stops the micro router
func (r *rtr) Stop() error {
	// stop the router
	if err := r.Stop(); err != nil {
		return fmt.Errorf("failed to stop router: %v", err)
	}

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

	// create new micro router and start advertising routes
	rtr := newRouter(service, r)

	log.Log("[router] starting to advertise")

	advertChan, err := rtr.Advertise()
	if err != nil {
		log.Logf("[router] failed to start: %s", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- rtr.PublishAdverts(advertChan)
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
