package server

import (
	"os"

	"github.com/urfave/cli/v2"
	"micro.dev/v4/service"
	"micro.dev/v4/service/client"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/router"
	"micro.dev/v4/service/server"
	"micro.dev/v4/service/server/grpc"
	"micro.dev/v4/util/muxer"
	"micro.dev/v4/util/proxy"
	grpcProxy "micro.dev/v4/util/proxy/grpc"
)

var (
	// name of the micro network
	network = "micro"
	// address is the network address
	address = ":8085"
	// netAddress is the rpc address
	netAddress = ":8443"

	// Flags specific to the network
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Set the address of the network service",
			EnvVars: []string{"MICRO_NETWORK_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "gateway",
			Usage:   "Set the default gateway",
			EnvVars: []string{"MICRO_NETWORK_GATEWAY"},
		},
		&cli.StringFlag{
			Name:    "network",
			Usage:   "Set the micro network name: micro",
			EnvVars: []string{"MICRO_NETWORK_NAME"},
		},
	}
)

// Run runs the micro server
func Run(ctx *cli.Context) error {
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}
	if len(ctx.String("network")) > 0 {
		network = ctx.String("network")
	}

	// Initialise the local service
	service := service.New(
		service.Name("network"),
		service.Address(address),
	)

	gateway := ctx.String("gateway")
	id := service.Server().Options().Id

	// increase the client retries
	client.DefaultClient.Init(
		client.Retries(3),
	)

	// local tunnel router
	rtr := router.DefaultRouter

	rtr.Init(
		router.Network(network),
		router.Id(id),
		router.Gateway(gateway),
		router.Cache(),
	)

	// local proxy using grpc
	// TODO: reenable after PR
	localProxy := grpcProxy.NewProxy(
		proxy.WithRouter(rtr),
		proxy.WithClient(service.Client()),
	)

	// local mux
	localMux := muxer.New("network", localProxy)

	// set the handler
	srv := grpc.NewServer(
		server.Name("network"),
		server.Address(netAddress),
		server.WithRouter(localMux),
	)

	// start the grpc server
	if err := srv.Start(); err != nil {
		log.Fatal("Error starting network: %v", err)
	}

	log.Infof("Network [%s] listening on %s", network, netAddress)

	if err := service.Run(); err != nil {
		log.Errorf("Network %s failed: %v", network, err)
		os.Exit(1)
	}

	// stop the grpc server
	return srv.Stop()
}
