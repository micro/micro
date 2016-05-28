package main

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/micro/micro/examples/template/srv/handler"
	"github.com/micro/micro/examples/template/srv/subscriber"

	example "github.com/micro/micro/examples/template/srv/proto/example"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.template"),
		micro.Version("latest"),
	)

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))


	// Register Struct as Subscriber
	service.Server().Subscribe(
		service.Server().NewSubscriber("topic.go.micro.srv.template", new(subscriber.Example)),
	)

	// Register Function as Subscriber
	service.Server().Subscribe(
		service.Server().NewSubscriber("topic.go.micro.srv.template", subscriber.Handler),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
