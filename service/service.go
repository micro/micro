// Package service provides a micro service
package service

import (
	"fmt"
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
	"github.com/micro/micro/v2/plugin"

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

type srvCommand struct {
	Name    string
	Command func(ctx *ccli.Context, srvOpts ...micro.Option)
	Flags   []ccli.Flag
}

var srvCommands = []srvCommand{
	{
		Name:    "auth",
		Command: auth.Run,
	},
	{
		Name:    "broker",
		Command: broker.Run,
	},
	{
		Name:    "config",
		Command: config.Run,
		Flags:   config.Flags,
	},
	{
		Name:    "debug",
		Command: debug.Run,
		Flags:   debug.Flags,
	},
	{
		Name:    "health",
		Command: health.Run,
		Flags:   health.Flags,
	},
	{
		Name:    "network",
		Command: network.Run,
		Flags:   network.Flags,
	},
	{
		Name:    "registry",
		Command: registry.Run,
	},
	{
		Name:    "router",
		Command: router.Run,
		Flags:   router.Flags,
	},
	{
		Name:    "runtime",
		Command: runtime.Run,
		Flags:   runtime.Flags,
	},
	{
		Name:    "store",
		Command: store.Run,
	},
	{
		Name:    "tunnel",
		Command: tunnel.Run,
		Flags:   tunnel.Flags,
	},
}

func Commands(options ...micro.Option) []*ccli.Command {
	// move newAction outside the loop and pass c as an arg to
	// set the scope of the variable
	newAction := func(c srvCommand) func(ctx *ccli.Context) error {
		return func(ctx *ccli.Context) error {
			c.Command(ctx, options...)
			return nil
		}
	}

	subcommands := make([]*ccli.Command, len(srvCommands))
	for i, c := range srvCommands {
		// construct the command
		command := &ccli.Command{
			Name:   c.Name,
			Flags:  c.Flags,
			Usage:  fmt.Sprintf("Run micro %v", c.Name),
			Action: newAction(c),
		}

		// setup the plugins
		for _, p := range plugin.Plugins(plugin.Module(c.Name)) {
			if cmds := p.Commands(); len(cmds) > 0 {
				command.Subcommands = append(command.Subcommands, cmds...)
			}

			if flags := p.Flags(); len(flags) > 0 {
				command.Flags = append(command.Flags, flags...)
			}
		}

		// set the command
		subcommands[i] = command
	}

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
		Subcommands: subcommands,
	}

	// register global plugins and flags
	for _, p := range plugin.Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []*ccli.Command{command}
}
