// Package proxy is a cli proxy
package proxy

import (
	"os"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	bmem "github.com/micro/go-micro/broker/memory"
	"github.com/micro/go-micro/client"
	mucli "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/proxy"
	"github.com/micro/go-micro/proxy/grpc"
	"github.com/micro/go-micro/proxy/http"
	"github.com/micro/go-micro/proxy/mucp"
	"github.com/micro/go-micro/registry"
	rmem "github.com/micro/go-micro/registry/memory"
	"github.com/micro/go-micro/router"
	rs "github.com/micro/go-micro/router/service"
	"github.com/micro/go-micro/server"
	sgrpc "github.com/micro/go-micro/server/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/util/mux"
)

var (
	// Name of the proxy
	Name = "go.micro.proxy"
	// The address of the proxy
	Address = ":8081"
	// the proxy protocol
	Protocol = "grpc"
	// The endpoint host to route to
	Endpoint string
)

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("proxy")

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("endpoint")) > 0 {
		Endpoint = ctx.String("endpoint")
	}
	if len(ctx.String("protocol")) > 0 {
		Protocol = ctx.String("protocol")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// service opts
	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(ctx.GlobalInt("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.GlobalInt("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// set the context
	var popts []proxy.Option

	// create new router
	var r router.Router

	routerName := ctx.String("router")
	routerAddr := ctx.String("router_address")

	ropts := []router.Option{
		router.Id(server.DefaultId),
		router.Client(client.DefaultClient),
		router.Address(routerAddr),
		router.Registry(registry.DefaultRegistry),
	}

	// check if we need to use the router service
	switch {
	case routerName == "go.micro.router":
		r = rs.NewRouter(ropts...)
	case len(routerAddr) > 0:
		r = rs.NewRouter(ropts...)
	default:
		r = router.NewRouter(ropts...)
	}

	// start the router
	if err := r.Start(); err != nil {
		log.Logf("Proxy error starting router: %s", err)
		os.Exit(1)
	}

	popts = append(popts, proxy.WithRouter(r))

	// new proxy
	var p proxy.Proxy
	var srv server.Server

	// set endpoint
	if len(Endpoint) > 0 {
		switch {
		case strings.HasPrefix(Endpoint, "grpc://"):
			ep := strings.TrimPrefix(Endpoint, "grpc://")
			popts = append(popts, proxy.WithEndpoint(ep))
			p = grpc.NewProxy(popts...)
		case strings.HasPrefix(Endpoint, "http://"):
			// TODO: strip prefix?
			popts = append(popts, proxy.WithEndpoint(Endpoint))
			p = http.NewProxy(popts...)
		default:
			// TODO: strip prefix?
			popts = append(popts, proxy.WithEndpoint(Endpoint))
			p = mucp.NewProxy(popts...)
		}
	}

	// set based on protocol
	if p == nil && len(Protocol) > 0 {
		switch Protocol {
		case "http":
			p = http.NewProxy(popts...)
			// TODO: http server
		case "mucp":
			popts = append(popts, proxy.WithClient(mucli.NewClient()))
			p = mucp.NewProxy(popts...)

			srv = server.NewServer(
				server.Address(Address),
				// reset registry to memory
				server.Registry(rmem.NewRegistry()),
				// reset broker to memory
				server.Broker(bmem.NewBroker()),
				// hande it the router
				server.WithRouter(p),
			)
		default:
			p = mucp.NewProxy(popts...)

			srv = sgrpc.NewServer(
				server.Address(Address),
				// reset registry to memory
				server.Registry(rmem.NewRegistry()),
				// reset broker to memory
				server.Broker(bmem.NewBroker()),
				// hande it the router
				server.WithRouter(p),
			)
		}
	}

	if len(Endpoint) > 0 {
		log.Logf("Proxy [%s] serving endpoint: %s", p.String(), Endpoint)
	} else {
		log.Logf("Proxy [%s] serving protocol: %s", p.String(), Protocol)
	}

	// new service
	service := micro.NewService(srvOpts...)

	// create a new proxy muxer which includes the debug handler
	muxer := mux.New(Name, p)

	// set the router
	service.Server().Init(
		server.WithRouter(muxer),
	)

	// Start the proxy server
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	// Run internal service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop the server
	if err := srv.Stop(); err != nil {
		log.Fatal(err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "proxy",
		Usage: "Run the service proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "router",
				Usage:  "Set the router to use e.g default, go.micro.router",
				EnvVar: "MICRO_ROUTER",
			},
			cli.StringFlag{
				Name:   "router_address",
				Usage:  "Set the router address",
				EnvVar: "MICRO_ROUTER_ADDRESS",
			},
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the proxy http address e.g 0.0.0.0:8081",
				EnvVar: "MICRO_PROXY_ADDRESS",
			},
			cli.StringFlag{
				Name:   "protocol",
				Usage:  "Set the protocol used for proxying e.g mucp, grpc, http",
				EnvVar: "MICRO_PROXY_PROTOCOL",
			},
			cli.StringFlag{
				Name:   "endpoint",
				Usage:  "Set the endpoint to route to e.g greeter or localhost:9090",
				EnvVar: "MICRO_PROXY_ENDPOINT",
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
