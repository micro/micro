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
	net "github.com/micro/go-micro/v2/network"
	res "github.com/micro/go-micro/v2/network/resolver"
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
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/service/network/handler"
)

var (
	// name of the network service
	name = "go.micro.network"
	// name of the micro network
	network = "go.micro"
	// address is the network address
	address = ":8085"
	// set the advertise address
	advertise = ""
	// resolver is the network resolver
	resolver = "registry"
	// the tunnel token
	token = "micro"

	// Flags specific to the network
	Flags = []cli.Flag{
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
	}
)

// Run runs the micro server
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}
	if len(ctx.String("advertise")) > 0 {
		advertise = ctx.String("advertise")
	}
	if len(ctx.String("network")) > 0 {
		network = ctx.String("network")
	}
	if len(ctx.String("token")) > 0 {
		token = ctx.String("token")
	}

	var nodes []string
	if len(ctx.String("nodes")) > 0 {
		nodes = strings.Split(ctx.String("nodes"), ",")
	}
	if len(ctx.String("resolver")) > 0 {
		resolver = ctx.String("resolver")
	}
	var r res.Resolver
	switch resolver {
	case "dns":
		r = &dns.Resolver{}
	case "http":
		r = &http.Resolver{}
	case "registry":
		r = &registry.Resolver{}
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
		micro.Name(name),
		micro.RegisterTTL(time.Duration(ctx.Int("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.Int("register_interval"))*time.Second),
	)

	// create a tunnel
	tunOpts := []tunnel.Option{
		tunnel.Address(address),
		tunnel.Token(token),
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
		router.Network(network),
		router.Id(id),
		router.Registry(service.Options().Registry),
		router.Advertise(strategy),
		router.Gateway(gateway),
	)

	// create new network
	n := net.NewNetwork(
		net.Id(id),
		net.Name(network),
		net.Address(address),
		net.Advertise(advertise),
		net.Nodes(nodes...),
		net.Tunnel(tun),
		net.Router(rtr),
		net.Resolver(r),
	)

	// local proxy
	prx := mucp.NewProxy(
		proxy.WithRouter(rtr),
		proxy.WithClient(service.Client()),
		proxy.WithLink("network", n.Client()),
	)

	// create a handler
	h := server.DefaultRouter.NewHandler(
		&handler.Network{Network: n},
	)

	// register the handler
	server.DefaultRouter.Handle(h)

	// create a new muxer
	mux := mux.New(name, prx)

	// init server
	service.Server().Init(
		server.WithRouter(mux),
	)

	// set network server to proxy
	n.Server().Init(server.WithRouter(mux))

	// connect network
	if err := n.Connect(); err != nil {
		log.Fatalf("Network failed to connect: %v", err)
	}

	// netClose hard exits if we have problems
	netClose := func(net net.Network) error {
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

	log.Infof("Network [%s] listening on %s", network, address)

	if err := service.Run(); err != nil {
		log.Errorf("Network %s failed: %v", network, err)
		netClose(n)
		os.Exit(1)
	}

	// close the network
	netClose(n)
}
