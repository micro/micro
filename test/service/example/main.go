package main

import (
	"fmt"

	"example/handler"
	pb "example/proto"
	"github.com/micro/micro/v3/service"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.example"),
	)

	// Initialise service
	srv.Init()

	// Register Handler
	pb.RegisterExampleHandler(srv.Server(), new(handler.Example))

	// Run service
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
