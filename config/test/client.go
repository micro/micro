package main

import (
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/mucp"
)

func main() {
	src := mucp.NewSource(
		mucp.ServiceName("go.micro.config"))
	if err := config.Load(src); err != nil {
		panic(err)
	}

	config.Get("a")
}
