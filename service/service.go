// Package service provides a micro service
package service

import (
	"os"
	"strings"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	prox "github.com/micro/go-micro/v2/proxy"
	"github.com/micro/go-micro/v2/proxy/grpc"
	"github.com/micro/go-micro/v2/proxy/http"
	"github.com/micro/go-micro/v2/proxy/mucp"
	rt "github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/micro/v2/internal/helper"

	// clients
	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/client/bot"
	"github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/client/proxy"
	"github.com/micro/micro/v2/client/web"

	// services
	"github.com/micro/micro/v2/service/auth"
	"github.com/micro/micro/v2/service/broker"
	"github.com/micro/micro/v2/service/config"
	"github.com/micro/micro/v2/service/debug"
	"github.com/micro/micro/v2/service/handler/exec"
	"github.com/micro/micro/v2/service/handler/file"
	"github.com/micro/micro/v2/service/health"
	"github.com/micro/micro/v2/service/network"
	"github.com/micro/micro/v2/service/registry"
	"github.com/micro/micro/v2/service/router"
	"github.com/micro/micro/v2/service/runtime"
	"github.com/micro/micro/v2/service/store"
	"github.com/micro/micro/v2/service/tunnel"
)

func Run(ctx *ccli.Context, opts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "service"}))

	name := ctx.String("name")
	address := ctx.String("address")
	endpoint := ctx.String("endpoint")

	metadata := make(map[string]string)
	for _, md := range ctx.StringSlice("metadata") {
		parts := strings.Split(md, "=")
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		val := strings.Join(parts[1:], "=")

		// set the key/val
		metadata[key] = val
	}

	if len(metadata) > 0 {
		opts = append(opts, micro.Metadata(metadata))
	}

	if len(name) > 0 {
		opts = append(opts, micro.Name(name))
	}

	if len(address) > 0 {
		opts = append(opts, micro.Address(address))
	}

	if len(endpoint) == 0 {
		endpoint = prox.DefaultEndpoint
	}

	var p prox.Proxy

	switch {
	case strings.HasPrefix(endpoint, "grpc"):
		endpoint = strings.TrimPrefix(endpoint, "grpc://")
		p = grpc.NewProxy(prox.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "http"):
		p = http.NewProxy(prox.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "file"):
		endpoint = strings.TrimPrefix(endpoint, "file://")
		p = file.NewProxy(prox.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "exec"):
		endpoint = strings.TrimPrefix(endpoint, "exec://")
		p = exec.NewProxy(prox.WithEndpoint(endpoint))
	default:
		p = mucp.NewProxy(prox.WithEndpoint(endpoint))
	}

	// run the service if asked to
	if ctx.Args().Len() > 0 {
		args := []rt.CreateOption{
			rt.WithCommand(ctx.Args().Slice()...),
			rt.WithOutput(os.Stdout),
		}

		// create new local runtime
		r := rt.NewRuntime()

		// start the runtime
		r.Start()

		// register the service
		r.Create(&rt.Service{
			Name: name,
		}, args...)

		// stop the runtime
		defer func() {
			r.Delete(&rt.Service{
				Name: name,
			})
			r.Stop()
		}()
	}

	log.Infof("Service [%s] Serving %s at endpoint %s\n", p.String(), name, endpoint)

	// new service
	service := micro.NewService(opts...)

	// create new muxer
	//	muxer := mux.New(name, p)

	// set the router
	service.Server().Init(
		server.WithRouter(p),
	)

	// run service
	service.Run()
}

func Commands(options ...micro.Option) []*ccli.Command {
	command := &ccli.Command{
		Name:  "service",
		Usage: "Run a micro service",
		Action: func(ctx *ccli.Context) error {
			Run(ctx, options...)
			return nil
		},
		Flags: []ccli.Flag{
			&ccli.StringFlag{
				Name:    "name",
				Usage:   "Name of the service",
				EnvVars: []string{"MICRO_SERVICE_NAME"},
				Value:   "service",
			},
			&ccli.StringFlag{
				Name:    "address",
				Usage:   "Address of the service",
				EnvVars: []string{"MICRO_SERVICE_ADDRESS"},
			},
			&ccli.StringFlag{
				Name:    "endpoint",
				Usage:   "The local service endpoint (Defaults to localhost:9090); {http, grpc, file, exec}://path-or-address e.g http://localhost:9090",
				EnvVars: []string{"MICRO_SERVICE_ENDPOINT"},
			},
			&ccli.StringSliceFlag{
				Name:    "metadata",
				Usage:   "Add metadata as key-value pairs describing the service e.g owner=john@example.com",
				EnvVars: []string{"MICRO_SERVICE_METADATA"},
			},
		},
		Subcommands: []*ccli.Command{
			{
				Name:  "api",
				Usage: "Run micro api",
				Action: func(ctx *ccli.Context) error {
					api.Run(ctx)
					return nil
				},
			},
			{
				Name:  "auth",
				Usage: "Run micro auth",
				Action: func(ctx *ccli.Context) error {
					auth.Run(ctx)
					return nil
				},
			},
			{
				Name:  "bot",
				Usage: "Run micro bot",
				Flags: bot.Flags,
				Action: func(ctx *ccli.Context) error {
					bot.Run(ctx)
					return nil
				},
			},
			{
				Name:  "broker",
				Usage: "Run micro broker",
				Action: func(ctx *ccli.Context) error {
					broker.Run(ctx)
					return nil
				},
			},
			{
				Name:  "cli",
				Usage: "Run micro cli",
				Action: func(ctx *ccli.Context) error {
					cli.Run(ctx)
					return nil
				},
			},
			{
				Name:  "config",
				Usage: "Run micro config",
				Flags: config.Flags,
				Action: func(ctx *ccli.Context) error {
					config.Run(ctx)
					return nil
				},
			},
			{
				Name:  "debug",
				Usage: "Run micro debug",
				Flags: debug.Flags,
				Action: func(ctx *ccli.Context) error {
					debug.Run(ctx)
					return nil
				},
			},
			{
				Name:  "health",
				Usage: "Run micro health",
				Flags: health.Flags,
				Action: func(ctx *ccli.Context) error {
					health.Run(ctx)
					return nil
				},
			},
			{
				Name:  "network",
				Usage: "Run micro network",
				Action: func(ctx *ccli.Context) error {
					if err := helper.UnexpectedSubcommand(ctx); err != nil {
						return err
					}
					network.Run(ctx, options...)
					return nil
				},
			},
			{
				Name:  "proxy",
				Usage: "Run micro proxy",
				Action: func(ctx *ccli.Context) error {
					proxy.Run(ctx, options...)
					return nil
				},
			},
			{
				Name:  "registry",
				Usage: "Run micro registry",
				Action: func(ctx *ccli.Context) error {
					if err := helper.UnexpectedSubcommand(ctx); err != nil {
						return err
					}
					registry.Run(ctx, options...)
					return nil
				},
			},
			{
				Name:  "router",
				Usage: "Run micro network router",
				Action: func(ctx *ccli.Context) error {
					router.Run(ctx, options...)
					return nil
				},
			},
			{
				Name:  "runtime",
				Usage: "Run micro runtime",
				Flags: runtime.Flags,
				Action: func(ctx *ccli.Context) error {
					if err := helper.UnexpectedSubcommand(ctx); err != nil {
						return err
					}
					runtime.Run(ctx, options...)
					return nil
				},
			},
			{
				Name:  "store",
				Usage: "Run micro store",
				Action: func(ctx *ccli.Context) error {
					if err := helper.UnexpectedSubcommand(ctx); err != nil {
						return err
					}
					store.Run(ctx, options...)
					return nil
				},
			},
			{
				Name:  "tunnel",
				Usage: "Run micro tunnel",
				Flags: tunnel.Flags,
				Action: func(ctx *ccli.Context) error {
					if err := helper.UnexpectedSubcommand(ctx); err != nil {
						return err
					}
					tunnel.Run(ctx, options...)
					return nil
				},
			},
			{
				Name:  "web",
				Usage: "Run micro web",
				Flags: web.Flags,
				Action: func(ctx *ccli.Context) error {
					if err := helper.UnexpectedSubcommand(ctx); err != nil {
						return err
					}
					web.Run(ctx, options...)
					return nil
				},
			},
		},
	}

	return []*ccli.Command{command}
}
