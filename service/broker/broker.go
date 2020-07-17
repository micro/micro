// Package broker is the micro broker
package broker

import (
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	pb "github.com/micro/go-micro/v2/broker/service/proto"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/service/broker/handler"
)

var (
	// Name of the broker
	Name = "go.micro.broker"
	// The address of the broker
	Address = ":8001"
)

func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "broker"}))

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// service opts
	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// set address
	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// new service
	service := micro.NewService(srvOpts...)

	// connect to the broker
	service.Options().Broker.Connect()

	// register the broker handler
	pb.RegisterBrokerHandler(service.Server(), &handler.Broker{
		// using the mdns broker
		Broker: service.Options().Broker,
	})

	// run the service
	service.Run()
}
