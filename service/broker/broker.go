// Package broker is the micro broker
package broker

import (
	"time"

	"github.com/micro/cli/v2"
	pb "github.com/micro/go-micro/v2/broker/service/proto"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/service"
	"github.com/micro/micro/v2/service/broker/handler"
)

var (
	name    = "go.micro.broker"
	address = ":8001"
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
	srv.Options().Broker.Connect()

	// register the broker handler
	pb.RegisterBrokerHandler(srv.Server(), &handler.Broker{
		Broker: srv.Options().Broker,
	})

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
	return nil
}
