package main

import (
	"dep-test-service/handler"
	"dep-test-service/subscriber"
	"fmt"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	dep "dep-test-service/proto/dep"
	dependency "dependency"
)

func main() {
	// New Service
	service := service.NewService(
		micro.Name("go.micro.service.dep"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	fmt.Println(dependency.Hello)

	// Register Handler
	dep.RegisterDepHandler(service.Server(), new(handler.Dep))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.dep", service.Server(), new(subscriber.Dep))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
