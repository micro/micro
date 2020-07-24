package main

import (
	"fmt"

	"github.com/micro/micro/v2/service"
	"github.com/micro/micro/v2/service/config"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.config-read"),
		service.Version("latest"),
	)
	srv.Init()

	// get a value
	c := config.DefaultConfig
	fmt.Println("Value of key.subkey: ", c.Get("key", "subkey").String(""))
}
