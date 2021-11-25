package server

import (
	"time"

	pb "github.com/micro/micro/v3/proto/broker"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/broker/handler"
	"github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"
)

var (
	name    = "broker"
	address = ":8003"
)

// Run the micro broker
func Run(ctx *cli.Context) error {
	srvOpts := []service.Option{
		service.Name(name),
		service.Address(address),
	}

	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterInterval(i*time.Second))
	}

	// new service
	srv := service.New(srvOpts...)

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
