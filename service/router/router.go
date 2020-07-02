package router

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/router"
	pb "github.com/micro/go-micro/v2/router/service/proto"
	"github.com/micro/micro/v2/service/router/handler"
)

var (
	// Name of the router microservice
	Name = "go.micro.router"
	// Address is the router microservice bind address
	Address = ":8084"
	// Network is the network name
	Network = router.DefaultNetwork
	// Topic is router adverts topic
	Topic = "go.micro.router.adverts"
)

// Sub processes router events
type sub struct {
	router router.Router
}

// Process processes router adverts
func (s *sub) Process(ctx context.Context, advert *pb.Advert) error {
	log.Debugf("received advert from: %s", advert.Id)
	if advert.Id == s.router.Options().Id {
		log.Debug("skipping advert")
		return nil
	}

	var events []*router.Event
	for _, event := range advert.Events {
		route := router.Route{
			Service: event.Route.Service,
			Address: event.Route.Address,
			Gateway: event.Route.Gateway,
			Network: event.Route.Network,
			Link:    event.Route.Link,
			Metric:  event.Route.Metric,
		}

		e := &router.Event{
			Type:      router.EventType(event.Type),
			Timestamp: time.Unix(0, advert.Timestamp),
			Route:     route,
		}

		events = append(events, e)
	}

	a := &router.Advert{
		Id:        advert.Id,
		Type:      router.AdvertType(advert.Type),
		Timestamp: time.Unix(0, advert.Timestamp),
		TTL:       time.Duration(advert.Ttl),
		Events:    events,
	}

	if err := s.router.Process(a); err != nil {
		return fmt.Errorf("failed processing advert: %s", err)
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
		log.Errorf("failed to subscribe to adverts: %s", err)
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
			log.Debugf("error publishing advert: %v", err)
			return fmt.Errorf("error publishing advert: %v", err)
		}
	}

	return nil
}

// Close the micro router
func (r *rtr) Close() error {
	// close the router
	if err := r.Router.Close(); err != nil {
		return fmt.Errorf("failed to close router: %v", err)
	}

	return nil
}

// Run runs the micro server
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "router"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("network")) > 0 {
		Network = ctx.String("network")
	}
	// default gateway address
	var gateway string
	if len(ctx.String("gateway")) > 0 {
		gateway = ctx.String("gateway")
	}

	// advertise the best routes
	strategy := router.AdvertiseLocal

	if a := ctx.String("advertise_strategy"); len(a) > 0 {
		switch a {
		case "all":
			strategy = router.AdvertiseAll
		case "best":
			strategy = router.AdvertiseBest
		case "local":
			strategy = router.AdvertiseLocal
		case "none":
			strategy = router.AdvertiseNone
		}
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
		micro.RegisterTTL(time.Duration(ctx.Int("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.Int("register_interval"))*time.Second),
	)

	r := router.NewRouter(
		router.Id(service.Server().Options().Id),
		router.Address(service.Server().Options().Id),
		router.Network(Network),
		router.Registry(service.Options().Registry),
		router.Gateway(gateway),
		router.Advertise(strategy),
	)

	// register router handler
	pb.RegisterRouterHandler(
		service.Server(),
		&handler.Router{
			Router: r,
		},
	)

	// register the table handler
	pb.RegisterTableHandler(
		service.Server(),
		&handler.Table{
			Router: r,
		},
	)

	// create new micro router and start advertising routes
	rtr := newRouter(service, r)

	log.Info("starting to advertise")

	advertChan, err := rtr.Advertise()
	if err != nil {
		log.Errorf("failed to advertise: %s", err)
		log.Info("attempting to stop the router")
		if err := rtr.Close(); err != nil {
			log.Errorf("failed to close: %s", err)
			os.Exit(1)
		}
		os.Exit(1)
	}

	// error channel to collect router errors
	errChan := make(chan error, 2)

	go func() {
		errChan <- rtr.PublishAdverts(advertChan)
	}()

	go func() {
		errChan <- service.Run()
	}()

	// we block here until either service or server fails
	if err := <-errChan; err != nil {
		log.Errorf("error running the router: %v", err)
	}

	log.Info("attempting to close the router")

	// close the router
	if err := r.Close(); err != nil {
		log.Errorf("failed to close: %s", err)
		os.Exit(1)
	}

	log.Info("successfully closed")
}

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "router",
		Usage: "Run the micro network router",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the micro router address :9093",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "network",
				Usage:   "Set the micro network name: local",
				EnvVars: []string{"MICRO_NETWORK_NAME"},
			},
			&cli.StringFlag{
				Name:    "gateway",
				Usage:   "Set the micro default gateway address. Defaults to none.",
				EnvVars: []string{"MICRO_GATEWAY_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "advertise_strategy",
				Usage:   "Set the advertise strategy; all, best, local, none",
				EnvVars: []string{"MICRO_ROUTER_ADVERTISE_STRATEGY"},
			},
		},
		Action: func(ctx *cli.Context) error {
			Run(ctx, options...)
			return nil
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

	return []*cli.Command{command}
}
