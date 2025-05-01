package main

import (
	"context"
	"fmt"

	hello "github.com/micro/micro/examples/greeter/srv/proto/hello"
	roundrobin "github.com/micro/plugins/v5/wrapper/select/roundrobin"
	"go-micro.dev/v5"
)

func main() {
	service := micro.NewService(
		micro.Name("client"),
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// parse command line flags
	service.Init()

	cli := hello.NewSayService("go.micro.srv.greeter", service.Client())

	response, err := cli.Hello(context.TODO(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v", response)
}
