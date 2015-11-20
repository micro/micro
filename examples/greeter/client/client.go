package main

import (
	"fmt"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	c "github.com/micro/go-micro/context"
	hello "github.com/micro/micro/examples/greeter/server/proto/hello"

	"golang.org/x/net/context"
)

func main() {
	cmd.Init()

	// Create new request to service go.micro.srv.greeter, method Say.Hello
	req := client.NewRequest("go.micro.srv.greeter", "Say.Hello", &hello.Request{
		Name: "John",
	})

	// Set arbitrary headers in context
	ctx := c.WithMetadata(context.Background(), map[string]string{
		"X-User-Id": "john",
		"X-From-Id": "script",
	})

	rsp := &hello.Response{}

	// Call service
	if err := client.Call(ctx, req, rsp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
