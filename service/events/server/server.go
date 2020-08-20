package server

import (
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v3/service"
	pb "github.com/micro/micro/v3/service/events/proto"
	"github.com/micro/micro/v3/service/logger"
)

// Run the micro broker
func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("events"),
	)

	// register the broker handler
	pb.RegisterStreamHandler(srv.Server(), new(evStream))

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	return nil
}
