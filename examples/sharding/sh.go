package main

import (
	"context"
	"fmt"

	hello "github.com/micro/micro/examples/greeter/srv/proto/hello"
	shard "github.com/micro/plugins/v5/wrapper/select/shard"
	"go-micro.dev/v5"
)

func main() {
	wrapper := shard.NewClientWrapper("X-From-User")

	service := micro.NewService(
		micro.Name("go.micro.api.greeter"),
		micro.WrapClient(wrapper),
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
