package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/micro/v2/test/example-service/handler"
	example "github.com/micro/micro/v2/test/example-service/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.example"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
