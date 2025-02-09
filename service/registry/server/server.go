package server

import (
	pb "github.com/micro/micro/v5/proto/registry"
	"github.com/micro/micro/v5/service"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/registry/handler"
	"github.com/urfave/cli/v2"
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
