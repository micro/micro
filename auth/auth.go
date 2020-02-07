package auth

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/micro/v2/auth/api"
	"github.com/micro/micro/v2/auth/handler"
)

var (
	// Name of the service
	Name = "go.micro.auth"
	// Address of the service
	Address = ":8010"
)

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("auth")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// setup service
	srvOpts = append(srvOpts, micro.Name(Name))
	service := micro.NewService(srvOpts...)

	// run service
	pb.RegisterAuthHandler(service.Server(), handler.New())
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Commands(srvOpts ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "auth",
		Usage: "Run the auth service",
		Action: func(ctx *cli.Context) error {
			run(ctx)
			return nil
		},

		Subcommands: append([]*cli.Command{
			{
				Name:        "api",
				Usage:       "Run the auth api",
				Description: "Run the auth api",
				Action: func(ctx *cli.Context) error {
					api.Run(ctx, srvOpts...)
					return nil
				},
			},
		}),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the auth http address e.g 0.0.0.0:8010",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
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
