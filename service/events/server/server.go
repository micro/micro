package server

import (
	pb "github.com/micro/micro/v5/proto/events"
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/events/handler"
	"github.com/micro/micro/v5/service/logger"
	"github.com/urfave/cli/v2"
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
