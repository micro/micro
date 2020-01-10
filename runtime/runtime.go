// Package runtime is the micro runtime
package runtime

import (
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	pb "github.com/micro/go-micro/runtime/service/proto"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/runtime/handler"
)

var (
	// Name of the runtime
	Name = "go.micro.runtime"
	// Address of the runtime
	Address = ":8088"
)

// Run the runtime service
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("runtime")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}

	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// create runtime
	muRuntime := *cmd.DefaultCmd.Options().Runtime

	// use default store
	muStore := *cmd.DefaultCmd.Options().Store

	// create a new runtime manager
	manager := newManager(muRuntime, muStore)

	// start the manager
	if err := manager.Start(); err != nil {
		log.Logf("failed to start: %s", err)
		os.Exit(1)
	}

	// append name
	srvOpts = append(srvOpts, micro.Name(Name))

	// new service
	service := micro.NewService(srvOpts...)

	// register the runtime handler
	pb.RegisterRuntimeHandler(service.Server(), &handler.Runtime{
		// Client to publish events
		Client: micro.NewEvent("go.micro.runtime.events", service.Client()),
		// using the micro runtime
		Runtime: manager,
	})

	// start runtime service
	if err := service.Run(); err != nil {
		log.Logf("error running service: %v", err)
	}

	// stop the manager
	if err := manager.Stop(); err != nil {
		log.Logf("failed to stop: %s", err)
		os.Exit(1)
	}
}

// Flags is shared flags so we don't have to continually re-add
func Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Set the name of the service to run",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to run",
			Value: "latest",
		},
		cli.StringFlag{
			Name:  "source",
			Usage: "Set the source url of the service e.g /path/to/source",
		},
		cli.BoolFlag{
			Name:  "local",
			Usage: "Set to run the service from local path",
		},
		cli.StringSliceFlag{
			Name:  "env",
			Usage: "Set the environment variables e.g. foo=bar",
		},
		cli.BoolFlag{
			Name:  "runtime",
			Usage: "Return the runtime services",
		},
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := []cli.Command{
		{
			Name:  "runtime",
			Usage: "Run the micro runtime",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Set the registry http address e.g 0.0.0.0:8088",
					EnvVar: "MICRO_SERVER_ADDRESS",
				},
				cli.StringSliceFlag{
					Name:   "env",
					Usage:  "Set environment variables for all services e.g. foo=bar",
					EnvVar: "MICRO_RUNTIME_ENV",
				},
			},
			Action: func(ctx *cli.Context) {
				Run(ctx, options...)
			},
		},
		{
			// In future we'll also have `micro run [x]` hence `micro run service` requiring "service"
			Name:  "run",
			Usage: "Run a service e.g micro run service version",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				runService(ctx, options...)
			},
		},
		{
			Name:  "kill",
			Usage: "Kill removes a running service e.g micro kill service",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				killService(ctx, options...)
			},
		},
		{
			Name:  "ps",
			Usage: "Ps returns status of a running service or lists all running services e.g. micro ps",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				getService(ctx, options...)
			},
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command[0].Subcommands = append(command[0].Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command[0].Flags = append(command[0].Flags, flags...)
		}
	}

	return command
}
