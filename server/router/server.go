package router

import (
	pb "github.com/micro/micro/v3/proto/router"
	"github.com/micro/micro/v3/service"
	muregistry "github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/router"
	"github.com/micro/micro/v3/service/router/registry"
	"github.com/urfave/cli/v2"
)

var (
	// name of the router microservice
	name = "router"
	// address is the router microservice bind address
	address = ":8084"
	// network is the network name
	network = router.DefaultNetwork

	// Flags specific to the router
	Flags = []cli.Flag{
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
	}
)

// Run the micro router
func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}
	if len(ctx.String("network")) > 0 {
		network = ctx.String("network")
	}
	// default gateway address
	var gateway string
	if len(ctx.String("gateway")) > 0 {
		gateway = ctx.String("gateway")
	}

	// Initialise service
	srv := service.New(
		service.Name(name),
		service.Address(address),
	)

	r := registry.NewRouter(
		router.Id(srv.Server().Options().Id),
		router.Address(srv.Server().Options().Id),
		router.Network(network),
		router.Registry(muregistry.DefaultRegistry),
		router.Gateway(gateway),
	)

	// register handlers
	pb.RegisterRouterHandler(srv.Server(), &Router{Router: r})
	pb.RegisterTableHandler(srv.Server(), &Table{Router: r})

	return srv.Run()
}
