package server

import (
	pb "github.com/micro/micro/v5/proto/broker"
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/broker"
	"github.com/micro/micro/v5/service/broker/handler"
	"github.com/micro/micro/v5/service/logger"
	"github.com/urfave/cli/v2"
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
