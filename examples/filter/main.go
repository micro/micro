package main

import (
	"context"
	"fmt"

	"github.com/micro/micro/examples/filter/version"
	proto "github.com/micro/micro/examples/service/proto"
	"go-micro.dev/v5"
)

func main() {
	service := micro.NewService()
	service.Init()

	greeter := proto.NewGreeterService("greeter", service.Client())

	rsp, err := greeter.Hello(
		// provide a context
		context.TODO(),
		// provide the request
		&proto.Request{Name: "John"},
		// set the filter
		version.Filter("latest"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Greeting)
}
