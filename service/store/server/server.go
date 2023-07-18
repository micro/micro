package store

import (
	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/store"
	"micro.dev/v4/service"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/store/handler"
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
