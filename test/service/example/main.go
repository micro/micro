package main

import (
	"example/handler"
	pb "example/proto"

	"github.com/micro/micro/v3/service"
)

func main() {
	// Register Handler
	pb.RegisterExampleHandler(new(handler.Example))

	// Run the service
	service.Run()
}
