package router

import (
	"os"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
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
	// Router is the router bind address
	Router = ":9093"
	// Network is the router network id
	Network = "local"
)

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
	if len(ctx.String("router")) > 0 {
		Router = ctx.String("router")
	}
	if len(ctx.String("network")) > 0 {
		Network = ctx.String("network")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	r := router.NewRouter(
		router.ID(service.Server().Options().Id),
		router.Address(Router),
		router.Network(Network),
		router.Registry(service.Client().Options().Registry),
	)

	// register router handler
	pb.RegisterRouterHandler(
		service.Server(),
		&handler.Router{Router: r},
	)

	// channel to collect errors
	errChan := make(chan error, 2)

	// WaitGroup to track goroutines
	var wg sync.WaitGroup

	// Start the micro server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Log("[router] starting micro router")
		errChan <- r.Advertise()
	}()

	// Start the micro server service
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

	// stop the server
	if err := r.Stop(); err != nil {
		log.Logf("[router] error stopping the router: %v", err)
		os.Exit(1)
	}

	// wait for all the goroutines to stop
	wg.Wait()

	log.Logf("[router] successfully stopped")
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "router",
		Usage: "Run the micro network router",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "router",
				Usage:  "Set the micro router address :9093",
				EnvVar: "MICRO_ROUTER_ADDRESS",
			},
			cli.StringFlag{
				Name:   "network",
				Usage:  "Set the micro network address: local",
				EnvVar: "MICRO_NETWORK_ADDRESS",
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
