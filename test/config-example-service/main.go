package main

import (
	"fmt"

	"github.com/micro/go-micro/v2"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.config-read"),
		micro.Version("latest"),
	)

	// create a new config
	c := service.Options().Config

	// set a value
	fmt.Println("Value of key.subkey: ", c.Get("key", "subkey").String(""))
}
