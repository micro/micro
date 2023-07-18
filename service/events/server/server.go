package server

import (
	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/events"
	"micro.dev/v4/service"
	"micro.dev/v4/service/events/handler"
	"micro.dev/v4/service/logger"
)

// Run the micro broker
func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("events"),
		service.Address(":8005"),
	)

	// register the handlers
	pb.RegisterStreamHandler(srv.Server(), new(handler.Stream))
	pb.RegisterStoreHandler(srv.Server(), new(handler.Store))

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	return nil
}
