package store

import (
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/store/memory"
	"github.com/micro/go-micro/store/service/handler"
	pb "github.com/micro/go-micro/store/service/proto"
	"github.com/micro/go-micro/util/log"
)

var (
	// Name of the tunnel service
	Name = "go.micro.store"
	// Address is the tunnel address
	Address = ":8002"
)

// run runs the micro server
func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("store")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	// TODO: allow flag flipping of backend store
	pb.RegisterStoreHandler(service.Server(), &handler.Store{
		Store: memory.NewStore(),
	})

	// start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "store",
		Usage: "Run the micro store service",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the micro tunnel address :8002",
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
