package server

import (
	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/broker"
	"micro.dev/v4/service"
	"micro.dev/v4/service/broker"
	"micro.dev/v4/service/broker/handler"
	"micro.dev/v4/service/logger"
)

var (
	address = ":8003"
)

// Run the micro broker
func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("broker"),
		service.Address(address),
	)

	// connect to the broker
	broker.DefaultBroker.Connect()

	// register the broker Broker
	pb.RegisterBrokerHandler(srv.Server(), new(handler.Broker))

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	return nil
}
