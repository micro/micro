package server

import (
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	pb "github.com/micro/go-micro/v2/store/service/proto"
	"github.com/micro/micro/v2/service"
	"github.com/micro/micro/v2/service/store"
)

var (
	// name of the store service
	name = "go.micro.store"
	// address is the store address
	address = ":8002"
)

// Run micro store
func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}

	// Initialise service
	service := service.New(
		service.Name(name),
		service.Address(address),
	)

	// the store handler
	pb.RegisterStoreHandler(service.Server(), &handler{
		store:  store.DefaultStore,
		stores: make(map[string]bool),
	})

	// start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
