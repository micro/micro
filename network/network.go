// package network implements micro network node
package network

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/network"
	"github.com/micro/go-micro/network/handler"
	"github.com/micro/go-micro/network/resolver"
	"github.com/micro/go-micro/network/resolver/dns"
	"github.com/micro/go-micro/network/resolver/http"
	"github.com/micro/go-micro/network/resolver/registry"
	"github.com/micro/go-micro/proxy"
	"github.com/micro/go-micro/proxy/mucp"
	"github.com/micro/go-micro/router"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/tunnel"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/util/mux"
	mcli "github.com/micro/micro/cli"
	"github.com/micro/micro/network/api"
)

var (
	// Name of the network service
	Name = "go.micro.network"
	// Name of the micro network
	Network = "go.micro"
	// Address is the network address
	Address = ":8085"
	// Set the advertise address
	Advertise = ""
	// Resolver is the network resolver
	Resolver = "registry"
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
	if len(ctx.String("advertise")) > 0 {
		Advertise = ctx.String("advertise")
	}
	if len(ctx.String("network")) > 0 {
		Network = ctx.String("network")
	}
	var nodes []string
	if len(ctx.String("nodes")) > 0 {
		nodes = strings.Split(ctx.String("nodes"), ",")
	}
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
	}
	var res resolver.Resolver
	switch Resolver {
	case "dns":
		res = &dns.Resolver{}
	case "http":
		res = &http.Resolver{}
	case "registry":
		res = &registry.Resolver{}
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	// create a tunnel
	tun := tunnel.NewTunnel(
		tunnel.Address(Address),
		tunnel.Nodes(nodes...),
	)

	// local tunnel router
	rtr := router.NewRouter(
		router.Network(Network),
		router.Id(service.Server().Options().Id),
		router.Registry(service.Client().Options().Registry),
	)

	// creaate new network
	net := network.NewNetwork(
		network.Id(service.Server().Options().Id),
		network.Name(Network),
		network.Address(Address),
		network.Advertise(Advertise),
		network.Peers(nodes...),
		network.Tunnel(tun),
		network.Router(rtr),
		network.Resolver(res),
	)

	// local proxy
	prx := mucp.NewProxy(
		proxy.WithRouter(rtr),
		proxy.WithClient(service.Client()),
		proxy.WithLink("network", net.Client()),
	)

	// create a handler
	h := server.DefaultRouter.NewHandler(
		&handler.Network{
			Network: net,
		},
	)

	// register the handler
	server.DefaultRouter.Handle(h)

	// create a new muxer
	mux := mux.New(Name, prx)

	// init server
	service.Server().Init(
		server.WithRouter(mux),
	)

	// set network server to proxy
	net.Server().Init(
		server.WithRouter(mux),
	)

	// connect network
	if err := net.Connect(); err != nil {
		log.Logf("Network failed to connect: %v", err)
		os.Exit(1)
	}

	// netClose hard exits if we have problems
	netClose := func(net network.Network) error {
		errChan := make(chan error, 1)

		go func() {
			errChan <- net.Close()
		}()

		select {
		case err := <-errChan:
			return err
		case <-time.After(time.Second):
			return errors.New("Network timeout closing")
		}
	}

	log.Logf("Network [%s] listening on %s", Name, Address)

	if err := service.Run(); err != nil {
		log.Logf("Network %s failed: %v", Name, err)
		netClose(net)
		os.Exit(1)
	}

	// close the network
	netClose(net)
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "network",
		Usage: "Run the micro network node",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the micro network address :8085",
				EnvVar: "MICRO_NETWORK_ADDRESS",
			},
			cli.StringFlag{
				Name:   "advertise",
				Usage:  "Set the micro network address to advertise",
				EnvVar: "MICRO_NETWORK_ADVERTISE",
			},
			cli.StringFlag{
				Name:   "network",
				Usage:  "Set the micro network name: go.micro",
				EnvVar: "MICRO_NETWORK",
			},
			cli.StringFlag{
				Name:   "nodes",
				Usage:  "Set the micro network nodes to connect to. This can be a comma separated list.",
				EnvVar: "MICRO_NETWORK_NODES",
			},
			cli.StringFlag{
				Name:   "resolver",
				Usage:  "Set the micro network resolver. This can be a comma separated list.",
				EnvVar: "MICRO_NETWORK_RESOLVER",
			},
		},
		Subcommands: append([]cli.Command{{
			Name:        "api",
			Description: "Run the network api",
			Action: func(ctx *cli.Context) {
				api.Run(ctx)
			},
		}}, mcli.NetworkCommands()...),
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
