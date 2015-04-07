package main

import (
	"fmt"

	"code.google.com/p/goprotobuf/proto"
	"github.com/asim/go-micro/client"
	hello "github.com/asim/micro/examples/greeter/server/proto/hello"
)

func main() {
	// Create new request to service go.micro.service.go-template, method Example.Call
	req := client.NewRequest("go.micro.srv.greeter", "Say.Hello", &hello.Request{
		Name: proto.String("John"),
	})

	// Set arbitrary headers
	req.Headers().Set("X-User-Id", "john")
	req.Headers().Set("X-From-Id", "script")

	rsp := &hello.Response{}

	// Call service
	if err := client.Call(req, rsp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.GetMsg())
}
