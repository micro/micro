package server

import (
	"log"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
)

var (
	// Name of the server
	Name = "go.micro.server"
	// Address to bind to
	Address = ":8083"
	// Network address to bind to
	Network = ":9093"
)

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("network")) > 0 {
		Network = ctx.String("network")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// Initi Mciro server
	var opts []server.Option
	opts = append(opts, server.Address(Address))

	// new server
	srv := server.NewServer()
	srv.Init(opts...)

	// service opts
	srvOpts = append(srvOpts, micro.Server(srv))
	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(ctx.GlobalInt("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.GlobalInt("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// Initialise service
	service := micro.NewService(srvOpts...)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "server",
		Usage: "Run the micro network server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "network",
				Usage:  "Set the micro network address 9093",
				EnvVar: "MICRO_NETWORK_ADDRESS",
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
