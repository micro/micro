// Package service provides a micro service
package service

import (
	"os"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/proxy"
	"github.com/micro/go-micro/v2/proxy/grpc"
	"github.com/micro/go-micro/v2/proxy/http"
	"github.com/micro/go-micro/v2/proxy/mucp"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/micro/v2/service/handler/exec"
	"github.com/micro/micro/v2/service/handler/file"
)

func run(ctx *cli.Context, opts ...micro.Option) {
	log.Name("service")

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
		endpoint = proxy.DefaultEndpoint
	}

	var p proxy.Proxy

	switch {
	case strings.HasPrefix(endpoint, "grpc"):
		endpoint = strings.TrimPrefix(endpoint, "grpc://")
		p = grpc.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "http"):
		p = http.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "file"):
		endpoint = strings.TrimPrefix(endpoint, "file://")
		p = file.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "exec"):
		endpoint = strings.TrimPrefix(endpoint, "exec://")
		p = exec.NewProxy(proxy.WithEndpoint(endpoint))
	default:
		p = mucp.NewProxy(proxy.WithEndpoint(endpoint))
	}

	// run the service if asked to
	if ctx.Args().Len() > 0 {
		args := []runtime.CreateOption{
			runtime.WithCommand(ctx.Args().Slice()...),
			runtime.WithOutput(os.Stdout),
		}

		// create new local runtime
		r := runtime.NewRuntime()

		// start the runtime
		r.Start()

		// register the service
		r.Create(&runtime.Service{
			Name: name,
		}, args...)

		// stop the runtime
		defer func() {
			r.Delete(&runtime.Service{
				Name: name,
			})
			r.Stop()
		}()
	}

	log.Logf("Service [%s] Serving %s at endpoint %s\n", p.String(), name, endpoint)

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

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "service",
		Usage: "Run a micro service",
		Action: func(ctx *cli.Context) error {
			run(ctx, options...)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Usage:   "Name of the service",
				EnvVars: []string{"MICRO_SERVICE_NAME"},
				Value:   "service",
			},
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Address of the service",
				EnvVars: []string{"MICRO_SERVICE_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "endpoint",
				Usage:   "The local service endpoint (Defaults to localhost:9090); {http, grpc, file, exec}://path-or-address e.g http://localhost:9090",
				EnvVars: []string{"MICRO_SERVICE_ENDPOINT"},
			},
			&cli.StringSliceFlag{
				Name:    "metadata",
				Usage:   "Add metadata as key-value pairs describing the service e.g owner=john@example.com",
				EnvVars: []string{"MICRO_SERVICE_METADATA"},
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

	return []*cli.Command{command}
}
