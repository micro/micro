// Package service provides a micro service
package service

import (
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/proxy"
	"github.com/micro/go-micro/proxy/grpc"
	"github.com/micro/go-micro/proxy/http"
	"github.com/micro/go-micro/proxy/mucp"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/util/mux"
)

func run(ctx *cli.Context, opts ...micro.Option) {
	log.Name("service")

	name := ctx.String("name")
	address := ctx.String("address")
	endpoint := ctx.String("endpoint")

	if len(name) > 0 {
		opts = append(opts, micro.Name(name))
	}

	if len(address) > 0 {
		opts = append(opts, micro.Address(address))
	}

	if len(endpoint) == 0 {
		endpoint = proxy.DefaultEndpoint
	}

	var p proxy.Proxy

	switch {
	case strings.HasPrefix(endpoint, "grpc"):
		p = grpc.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "http"):
		p = http.NewProxy(proxy.WithEndpoint(endpoint))
	default:
		p = mucp.NewProxy(proxy.WithEndpoint(endpoint))
	}

	log.Logf("Service [%s] Serving %s at endpoint %s\n", p.String(), name, endpoint)

	// new service
	service := micro.NewService(opts...)

	// create new muxer
	muxer := mux.New(name, p)

	// set the router
	service.Server().Init(
		server.WithRouter(muxer),
	)

	// run service
	service.Run()
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "service",
		Usage: "Run a micro service",
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "name",
				Usage:  "Name of the service",
				EnvVar: "MICRO_SERVICE_NAME",
			},
			cli.StringFlag{
				Name:   "address",
				Usage:  "Address of the service",
				EnvVar: "MICRO_SERVICE_ADDRESS",
			},
			cli.StringFlag{
				Name:   "endpoint",
				Usage:  "The local service endpoint. Defaults to localhost:9090",
				EnvVar: "MICRO_SERVICE_ENDPOINT",
			},
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
