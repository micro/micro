// Package debug allows to debug services
package debug

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/debug/handler"
	pb "github.com/micro/go-micro/debug/proto"
	"github.com/micro/go-micro/util/log"
)

var (
	// Name of the service
	Name = "go.micro.debug"
	// Address of the service
	Address = ":8089"
)

func getLogs(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("debug")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}
}

func getStats(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("debug")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}
}

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("debug")

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

	// append name
	srvOpts = append(srvOpts, micro.Name(Name))

	// new service
	service := micro.NewService(srvOpts...)

	pb.RegisterDebugHandler(service.Server(),
		handler.DefaultHandler,
	)

	// start debug service
	if err := service.Run(); err != nil {
		log.Logf("error running service: %v", err)
	}

	log.Logf("successfully stopped")
}

// Flags is shared flags so we don't have to continually re-add
func Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Set the name of the service to debug",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to debug",
			Value: "latest",
		},
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := []cli.Command{
		{
			Name:  "debug",
			Usage: "Run the micro debug service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Set the registry http address e.g 0.0.0.0:8089",
					EnvVar: "MICRO_SERVER_ADDRESS",
				},
			},
			Action: func(ctx *cli.Context) {
				run(ctx, options...)
			},
		},
		{
			Name:  "logs",
			Usage: "Get logs for a service",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				getLogs(ctx, options...)
			},
		},
		{
			Name:  "stats",
			Usage: "Get stats for a service",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				getStats(ctx, options...)
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
