package server

import (
	"os"

	"github.com/urfave/cli/v2"
	"micro.dev/v4/service"
	"micro.dev/v4/service/client"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/proxy"
	grpcProxy "micro.dev/v4/service/proxy/grpc"
	"micro.dev/v4/service/router"
	"micro.dev/v4/service/server"
	"micro.dev/v4/util/muxer"
)

var (
	// name of the network service
	name = "network"
	// name of the micro network
	networkName = "micro"
	// address is the network address
	address = ":8443"
	// peerAddress is the address the network peers on
	peerAddress = ":8085"

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
			EnvVars: []string{"MICRO_NETWORK"},
		},
	}
)

// Run runs the micro server
func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}
	if len(ctx.String("peer_address")) > 0 {
		peerAddress = ctx.String("peer_address")
	}
	if len(ctx.String("network")) > 0 {
		networkName = ctx.String("network")
	}

	// Initialise the local service
	service := service.New(
		service.Name(name),
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
		router.Network(networkName),
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
	localMux := muxer.New(name, localProxy)

	// init the local grpc server
	service.Server().Init(
		server.WithRouter(localMux),
	)

	log.Infof("Network [%s] listening on %s", networkName, peerAddress)

	if err := service.Run(); err != nil {
		log.Errorf("Network %s failed: %v", networkName, err)
		os.Exit(1)
	}

	return nil
}
