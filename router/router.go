package router

import (
	"os"
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
	// Router is the router gossip bind address
	Router = ":9093"
	// Network is the network id
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

	if _, err := r.Advertise(); err != nil {
		log.Logf("[router] failed to start: %s", err)
		os.Exit(1)
	}

	if err := service.Run(); err != nil {
		log.Logf("[router] failed with error %s", err)
		// TODO: we should probably stop the router here before bailing
		os.Exit(1)
	}

	// stop the router
	if err := r.Stop(); err != nil {
		log.Logf("[router] failed to stop: %s", err)
		os.Exit(1)
	}

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
