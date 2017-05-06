// Package run is a micro service runtime
package run

import (
	"fmt"
	"log"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/go-run/runtime"
	"github.com/micro/go-run/runtime/go"
	proto "github.com/micro/micro/run/proto"
)

var (
	Name = "go.micro.run"
)

func run(ctx *cli.Context) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	r := gorun.NewRuntime()

	// micro run github.com/my/service
	// args: github.com/my/service
	if len(ctx.Args()) > 0 {
		// look for flag to loop run
		if err := runtime.Run(r, ctx.Args().First()); err != nil {
			fmt.Println(err)
		}
		return
	}

	// Initialise Server
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	proto.RegisterRuntimeHandler(service.Server(), &handler{r})

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Commands() []cli.Command {
	command := cli.Command{
		Name:   "run",
		Usage:  "Run the micro runtime",
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

	return []cli.Command{command}
}
