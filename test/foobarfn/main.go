package main

import (
  log	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"foobarfn/handler"
	"foobarfn/subscriber"
)

func main() {
	// New Service
	function := micro.NewFunction(
		micro.Name("go.micro.function.foobarfn"),
		micro.Version("latest"),
	)

	// Initialise function
	function.Init()

	// Register Handler
	function.Handle(new(handler.Foobarfn))

	// Register Struct as Subscriber
	function.Subscribe("go.micro.function.foobarfn", new(subscriber.Foobarfn))

	// Run service
	if err := function.Run(); err != nil {
		log.Fatal(err)
	}
}
