package main

import (
	"dep-test-service/handler"
	"dep-test-service/subscriber"
	"fmt"

	dep "dep-test-service/proto/dep"
	dependency "dependency"

	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("service.dep"),
		service.Version("latest"),
	)

	// Initialise service
	srv.Init()
	fmt.Println(dependency.Hello)

	// Register Handler
	dep.RegisterDepHandler(srv.Server(), new(handler.Dep))

	// Register Struct as Subscriber
	service.RegisterSubscriber("service.dep", srv.Server(), new(subscriber.Dep))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
