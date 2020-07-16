package store

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	pb "github.com/micro/go-micro/v2/store/service/proto"
	mcli "github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/service/store/handler"
)

var (
	// Name of the store service
	Name = "go.micro.store"
	// Address is the store address
	Address = ":8002"
)

// Run runs the micro server
func Run(ctx *cli.Context, srvOpts ...micro.Option) error {
	log.Init(log.WithFields(map[string]interface{}{"service": "store"}))

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
	)

	// the store handler
	h := handler.New(service.Options().Store)
	pb.RegisterStoreHandler(service.Server(), h)

	// start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Commands is the cli interface for the store service
func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:        "store",
		Subcommands: mcli.StoreCommands(),
	}

	return []*cli.Command{command}
}
