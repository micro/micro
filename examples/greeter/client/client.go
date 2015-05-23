package main

import (
	"fmt"

	"github.com/myodc/go-micro/client"
	c "github.com/myodc/go-micro/context"
	hello "github.com/myodc/micro/examples/greeter/server/proto/hello"

	"golang.org/x/net/context"
)

func main() {
	// Create new request to service go.micro.service.go-template, method Example.Call
	req := client.NewRequest("go.micro.srv.greeter", "Say.Hello", &hello.Request{
		Name: "John",
	})

	// Set arbitrary headers in context
	ctx := c.WithMetaData(context.Background(), map[string]string{
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
