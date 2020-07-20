// Package broker is the micro broker
package broker

import (
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	pb "github.com/micro/go-micro/v2/broker/service/proto"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/service/broker/handler"
)

var (
	name    = "go.micro.broker"
	address = ":8001"
)

// Run the micro broker
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	srvOpts = append([]micro.Option{
		micro.Name(name),
		micro.Address(address),
	}, srvOpts...)

	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// new service
	srv := micro.NewService(srvOpts...)

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
}
