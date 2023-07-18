package server

import (
	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/registry"
	"micro.dev/v4/service"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/registry/handler"
)

var (
	// name of the registry
	name = "registry"
	// address of the registry
	address = ":8000"
	// topic to publish registry events to
	topic = "registry.events"
)

func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("registry"),
		service.Address(address),
	)
	// register the handler
	pb.RegisterRegistryHandler(srv.Server(), &handler.Registry{
		Event: service.NewEvent(topic),
	})

	// run the service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
