package main

import (
	"log"
	"time"

	"context"
	hello "github.com/micro/micro/examples/greeter/srv/proto/hello"
	"go-micro.dev/v5"
	"go-micro.dev/v5/transport"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		// setup a new transport with secure option
		micro.Transport(
			// create new transport
			transport.NewHTTPTransport(
				// set to automatically secure
				transport.Secure(true),
			),
		),
	)

	service.Init()

	hello.RegisterSayHandler(service.Server(), new(Say))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
