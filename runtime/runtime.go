// Package runtime is the micro runtime
package runtime

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/runtime/service/handler"
	pb "github.com/micro/go-micro/runtime/service/proto"
	"github.com/micro/go-micro/util/log"
)

var (
	// Name of the runtime
	Name = "go.micro.runtime"
	// Address of the runtime
	Address = ":8088"
)

func run(ctx *cli.Context, srvOpts ...micro.Option) {
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

	// append name
	srvOpts = append(srvOpts, micro.Name(Name))

	// new service
	service := micro.NewService(srvOpts...)

	// register the runtime handler
	pb.RegisterRuntimeHandler(service.Server(), &handler.Runtime{
		// using the mdns runtime
		Runtime: *cmd.DefaultCmd.Options().Runtime,
	})

	// run the service
	service.Run()
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "runtime",
		Usage: "Run the micro runtime",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the registry http address e.g 0.0.0.0:8000",
				EnvVar: "MICRO_SERVER_ADDRESS",
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
