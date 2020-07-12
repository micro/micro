package main

import (
	"fmt"

	"dep-test-service/handler"
	"dep-test-service/subscriber"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	dep "dep-test-service/proto/dep"
	dependency "dependency"
)

func main() {
	// New Service
	service := service.NewService(
		service.Name("go.micro.service.dep"),
		service.Version("latest"),
	)

	// Initialise service
	service.Init()
	fmt.Println(dependency.Hello)

	// Register Handler
	dep.RegisterDepHandler(service.Server(), new(handler.Dep))

	// Register Struct as Subscriber
	server := service.Server()
	server.Subscribe(
		server.NewSubscriber("go.micro.service.dep", new(subscriber.Dep)),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
