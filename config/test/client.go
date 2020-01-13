package main

import (
	"github.com/micro/go-micro/config"
	"github.com/micro/micro/config/source"
)

func main() {
	src := mucp.NewSource(
		mucp.Id("NAMESPACE:CONFIG"),
		mucp.ServiceName("go.micro.srv.config"))
	if err := config.Load(src); err != nil {
		panic(err)
	}

	config.Get("a")
}
