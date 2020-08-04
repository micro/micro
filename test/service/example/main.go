package main

import (
	"example/handler"
	pb "example/proto"

	"github.com/micro/micro/v3/service"
)

func main() {
	srv := service.New(service.Name("example"))

	// Register Handler
	pb.RegisterExampleHandler(srv.Server(), new(handler.Example))

	// Run the service
	srv.Run()
}
