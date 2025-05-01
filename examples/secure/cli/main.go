package main

import (
	"fmt"

	"go-micro.dev/v5/client"
	"go-micro.dev/v5/transport"

	hello "github.com/micro/micro/examples/greeter/srv/proto/hello"

	"context"
)

func init() {
	client.DefaultClient.Init(
		client.Transport(
			transport.NewHTTPTransport(transport.Secure(true)),
		),
	)
}

func main() {
	cl := hello.NewSayService("go.micro.srv.greeter", client.DefaultClient)

	rsp, err := cl.Hello(context.TODO(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
