package main

import (
	"context"
	"fmt"

	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
)

func main() {
	service := micro.NewService()
	service.Init()
	c := service.Client()

	request := c.NewRequest("greeter", "Greeter.Hello", "john", client.WithContentType("application/json"))
	var response string

	if err := c.Call(context.TODO(), request, &response); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
}
