package main

import (
	"dep-test-service/handler"
	"dep-test-service/subscriber"
	"fmt"

	dep "dep-test-service/proto/dep"
	dependency "dependency"

	"micro.dev/v4/service"
	log "micro.dev/v4/service/logger"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("dep"),
		service.Version("latest"),
	)

	// Initialise service
	srv.Init()
	fmt.Println(dependency.Hello)

	// Register Handler
	dep.RegisterDepHandler(srv.Server(), new(handler.Dep))

	// Register Struct as Subscriber
	service.Subscribe("dep", new(subscriber.Dep))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
