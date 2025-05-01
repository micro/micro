package main

import (
	"github.com/micro/micro/examples/auth/handler"
	pb "github.com/micro/micro/examples/auth/proto"

	"go-micro.dev/v5"
	"go-micro.dev/v5/logger"

	// Import the JWT auth plugin
	_ "github.com/micro/plugins/v5/auth/jwt"
)

var (
	name    = "helloworld"
	version = "latest"
)

func main() {

	// Create service
	srv := micro.NewService()

	srv.Init(
		micro.Name(name),
		micro.Version(version),
		micro.WrapHandler(NewAuthWrapper(srv)),
	)

	// Register handler
	if err := pb.RegisterHelloworldHandler(srv.Server(), new(handler.Helloworld)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

}
