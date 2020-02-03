package auth

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/micro/auth/handler"
)

func run(ctx *cli.Context) error {
	log.Name("auth")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// setup service
	service := micro.NewService(
		micro.Name("go.micro.auth"),
	)

	// run service
	pb.RegisterAuthHandler(service.Server(), handler.New())
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func Commands() []*cli.Command {
	command := &cli.Command{
		Name:   "auth",
		Usage:  "Run the auth service",
		Action: run,
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
