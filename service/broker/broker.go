// Package broker is the micro broker
package broker

import (
	"time"

	"github.com/micro/cli/v2"
	pb "github.com/micro/go-micro/v2/broker/service/proto"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/service/broker/handler"
)

var (
	// Name of the broker
	Name = "go.micro.broker"
	// The address of the broker
	Address = ":8001"
)

func init() {
	// register the commands
	cmd.Commands = append(app.Commands, Commands()...)
}

func Run(ctx *cli.Context, srvOpts ...service.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "broker"}))

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// service opts
	srvOpts = append(srvOpts, service.Name(Name))
	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterInterval(i*time.Second))
	}

	// set address
	if len(Address) > 0 {
		srvOpts = append(srvOpts, service.Address(Address))
	}

	// new service
	service := service.NewService(srvOpts...)

	// connect to the broker
	service.Options().Broker.Connect()

	// register the broker handler
	pb.RegisterBrokerHandler(service.Server(), &handler.Broker{
		// using the mdns broker
		Broker: service.Options().Broker,
	})

	// run the service
	service.Run()
}

func Commands(options ...service.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "broker",
		Usage: "Run the message broker",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the broker http address e.g 0.0.0.0:8001",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			},
		},
		Action: func(ctx *cli.Context) error {
			Run(ctx, options...)
			return nil
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
