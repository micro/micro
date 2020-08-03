package main

import (
	"example/handler"
	pb "example/proto"
)

func main() {
	// Register Handler
	pb.RegisterExampleService(new(handler.Example))

	// Run the service
	pb.RunExampleService()
}
