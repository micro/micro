package server

import (
	"strings"

	"github.com/urfave/cli/v2"
	"micro.dev/v4/service"
	bmem "micro.dev/v4/service/broker/memory"
	muclient "micro.dev/v4/service/client"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/proxy"
	"micro.dev/v4/service/proxy/grpc"
	"micro.dev/v4/service/proxy/http"
	"micro.dev/v4/service/registry/noop"
	murouter "micro.dev/v4/service/router"
	"micro.dev/v4/service/server"
	sgrpc "micro.dev/v4/service/server/grpc"
	"micro.dev/v4/util/muxer"
)

var (
	// Name of the proxy
	Name = "proxy"
	// The address of the proxy
	Address = ":8081"
	// Is gRPCWeb enabled
	GRPCWebEnabled = false
	// The address of the proxy
	GRPCWebAddress = ":8082"
	// the proxy protocol
	Protocol = "grpc"
	// The endpoint host to route to
	Endpoint string
)

func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
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
	// new service
	service := service.New(service.Name(Name))

	// set the context
	popts := []proxy.Option{
		proxy.WithRouter(murouter.DefaultRouter),
		proxy.WithClient(muclient.DefaultClient),
	}

	// set endpoint
	if len(Endpoint) > 0 {
		ep := Endpoint

		switch {
		case strings.HasPrefix(Endpoint, "grpc://"):
			ep = strings.TrimPrefix(Endpoint, "grpc://")
			Protocol = "grpc"
		case strings.HasPrefix(Endpoint, "http://"):
			Protocol = "http"
		}

		popts = append(popts, proxy.WithEndpoint(ep))
	}

	serverOpts := []server.Option{
		server.Name(Name),
		server.Address(Address),
		server.Registry(noop.NewRegistry()),
		server.Broker(bmem.NewBroker()),
	}

	// new proxy
	var p proxy.Proxy

	// set proxy
	switch Protocol {
	case "http":
		p = http.NewProxy(popts...)
		// TODO: http server
	default:
		// default to the grpc proxy
		p = grpc.NewProxy(popts...)
	}

	// wrap the proxy using the proxy's authHandler
	authOpt := server.WrapHandler(authHandler())
	serverOpts = append(serverOpts, authOpt)
	serverOpts = append(serverOpts, server.WithRouter(p))

	if len(Endpoint) > 0 {
		log.Infof("Proxy [%s] serving endpoint: %s", p.String(), Endpoint)
	} else {
		log.Infof("Proxy [%s] serving protocol: %s", p.String(), Protocol)
	}

	// create a new grpc server
	srv := sgrpc.NewServer(serverOpts...)

	// create a new proxy muxer which includes the debug handler
	muxer := muxer.New(Name, p)

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

	return nil
}

var (
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Set the proxy http address e.g 0.0.0.0:8081",
			EnvVars: []string{"MICRO_PROXY_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "protocol",
			Usage:   "Set the protocol used for proxying e.g mucp, grpc, http",
			EnvVars: []string{"MICRO_PROXY_PROTOCOL"},
		},
		&cli.StringFlag{
			Name:    "endpoint",
			Usage:   "Set the endpoint to route to e.g greeter or localhost:9090",
			EnvVars: []string{"MICRO_PROXY_ENDPOINT"},
		},
	}
)
