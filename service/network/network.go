// package network implements micro network node
package network

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/network"
	"github.com/micro/go-micro/v2/network/resolver"
	"github.com/micro/go-micro/v2/network/resolver/dns"
	"github.com/micro/go-micro/v2/network/resolver/http"
	"github.com/micro/go-micro/v2/network/resolver/registry"
	"github.com/micro/go-micro/v2/proxy"
	"github.com/micro/go-micro/v2/proxy/mucp"
	"github.com/micro/go-micro/v2/router"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/transport"
	"github.com/micro/go-micro/v2/transport/quic"
	"github.com/micro/go-micro/v2/tunnel"
	"github.com/micro/go-micro/v2/util/mux"
	mcli "github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/service/network/api"
	netdns "github.com/micro/micro/v2/service/network/dns"
	"github.com/micro/micro/v2/service/network/handler"
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
	// The tunnel token
	Token = "micro"
)

// Run runs the micro server
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "network"}))

	if ctx.Args().Len() > 0 {
		cli.ShowSubcommandHelp(ctx)
		os.Exit(1)
	}

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
	if len(ctx.String("advertise")) > 0 {
		Advertise = ctx.String("advertise")
	}
	if len(ctx.String("network")) > 0 {
		Network = ctx.String("network")
	}
	if len(ctx.String("token")) > 0 {
		Token = ctx.String("token")
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
		micro.RegisterTTL(time.Duration(ctx.Int("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.Int("register_interval"))*time.Second),
	)

	// create a tunnel
	tunOpts := []tunnel.Option{
		tunnel.Address(Address),
		tunnel.Token(Token),
	}

	if ctx.Bool("enable_tls") {
		config, err := helper.TLSConfig(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		config.InsecureSkipVerify = true

		tunOpts = append(tunOpts, tunnel.Transport(
			quic.NewTransport(transport.TLSConfig(config)),
		))
	}

	gateway := ctx.String("gateway")
	tun := tunnel.NewTunnel(tunOpts...)
	id := service.Server().Options().Id

	// local tunnel router
	rtr := router.NewRouter(
		router.Network(Network),
		router.Id(id),
		router.Registry(service.Client().Options().Registry),
		router.Advertise(strategy),
		router.Gateway(gateway),
	)

	// create new network
	net := network.NewNetwork(
		network.Id(id),
		network.Name(Network),
		network.Address(Address),
		network.Advertise(Advertise),
		network.Nodes(nodes...),
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
		log.Errorf("Network failed to connect: %v", err)
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

	log.Infof("Network [%s] listening on %s", Network, Address)

	if err := service.Run(); err != nil {
		log.Errorf("Network %s failed: %v", Network, err)
		netClose(net)
		os.Exit(1)
	}

	// close the network
	netClose(net)
}

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "network",
		Usage: "Run the micro network node",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the micro network address :8085",
				EnvVars: []string{"MICRO_NETWORK_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "advertise",
				Usage:   "Set the micro network address to advertise",
				EnvVars: []string{"MICRO_NETWORK_ADVERTISE"},
			},
			&cli.StringFlag{
				Name:    "gateway",
				Usage:   "Set the default gateway",
				EnvVars: []string{"MICRO_NETWORK_GATEWAY"},
			},
			&cli.StringFlag{
				Name:    "network",
				Usage:   "Set the micro network name: go.micro",
				EnvVars: []string{"MICRO_NETWORK"},
			},
			&cli.StringFlag{
				Name:    "nodes",
				Usage:   "Set the micro network nodes to connect to. This can be a comma separated list.",
				EnvVars: []string{"MICRO_NETWORK_NODES"},
			},
			&cli.StringFlag{
				Name:    "token",
				Usage:   "Set the micro network token for authentication",
				EnvVars: []string{"MICRO_NETWORK_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "resolver",
				Usage:   "Set the micro network resolver. This can be a comma separated list.",
				EnvVars: []string{"MICRO_NETWORK_RESOLVER"},
			},
			&cli.StringFlag{
				Name:    "advertise_strategy",
				Usage:   "Set the route advertise strategy; all, best, local, none",
				EnvVars: []string{"MICRO_NETWORK_ADVERTISE_STRATEGY"},
			},
		},
		Subcommands: append([]*cli.Command{
			{
				Name:        "api",
				Usage:       "Run the network api",
				Description: "Run the network api",
				Action: func(ctx *cli.Context) error {
					api.Run(ctx)
					return nil
				},
			},
			{
				Name:        "dns",
				Usage:       "Start a DNS resolver service that registers core nodes in DNS",
				Description: "Start a DNS resolver service that registers core nodes in DNS",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "provider",
						Usage:   "The DNS provider to use. Currently, only cloudflare is implemented",
						EnvVars: []string{"MICRO_NETWORK_DNS_PROVIDER"},
						Value:   "cloudflare",
					},
					&cli.StringFlag{
						Name:    "api-token",
						Usage:   "The provider's API Token.",
						EnvVars: []string{"MICRO_NETWORK_DNS_API_TOKEN"},
					},
					&cli.StringFlag{
						Name:    "zone-id",
						Usage:   "The provider's Zone ID.",
						EnvVars: []string{"MICRO_NETWORK_DNS_ZONE_ID"},
					},
					&cli.StringFlag{
						Name:    "token",
						Usage:   "Shared secret that must be presented to the service to authorize requests.",
						EnvVars: []string{"MICRO_NETWORK_DNS_TOKEN"},
					},
				},
				Action: func(ctx *cli.Context) error {
					netdns.Run(ctx)
					return nil
				},
				Subcommands: mcli.NetworkDNSCommands(),
			},
		}, mcli.NetworkCommands()...),
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
