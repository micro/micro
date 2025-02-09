package store

import (
	pb "github.com/micro/micro/v5/proto/store"
	"github.com/micro/micro/v5/service"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store/handler"
	"github.com/urfave/cli/v2"
)

var (
	// address is the store address
	address = ":8002"
)

// Run micro store
func Run(ctx *cli.Context) error {
	// Initialise service
	service := service.New(
		service.Name("store"),
		service.Address(address),
	)

	// the store handler
	pb.RegisterStoreHandler(service.Server(), &handler.Store{
		Stores: make(map[string]bool),
	})

	// the blob store handler
	pb.RegisterBlobStoreHandler(service.Server(), new(handler.BlobStore))

	// start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
