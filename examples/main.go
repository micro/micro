package main

import (
	"context"
	"io"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	pb "github.com/micro/go-micro/network/router/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.router"),
	)

	// Initialise service
	service.Init()

	client := pb.NewRouterService("go.micro.router", service.Client())

	id := "1"

	stream, err := client.Watch(context.Background(), &pb.WatchRequest{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			log.Logf("event stream disconnected")
			break
		}

		if err != nil {
			log.Logf("error receiving table event stream: %s", id)
			break
		}

		log.Logf("received event: %v", event)
	}
}
