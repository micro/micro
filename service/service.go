// Package service provides a micro service
package service

import (
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-proxy/router/http"
	"github.com/micro/go-proxy/router/mucp"
)

func run(ctx *cli.Context, opts ...micro.Option) {
	name := ctx.String("name")
	address := ctx.String("address")
	endpoint := ctx.String("endpoint")

	if len(name) > 0 {
		opts = append(opts, micro.Name(name))
	} else {
		name = server.DefaultName
	}

	if len(address) > 0 {
		opts = append(opts, micro.Address(address))
	}

	switch {
	case strings.HasPrefix(endpoint, "http"):
		opts = append(opts, http.WithRouter(&http.Router{
			Backend: endpoint,
		}))
	default:
		opts = append(opts, mucp.WithRouter(&mucp.Router{
			Name:    name,
			Backend: endpoint,
		}))
	}

	// new service
	service := micro.NewService(opts...)

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
				Usage:  "The local service endpoint",
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
