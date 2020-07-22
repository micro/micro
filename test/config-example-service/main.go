package main

import (
	"fmt"

	"github.com/micro/micro/v2/service"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.config-read"),
		service.Version("latest"),
	)
	srv.Init()

	// create a new config
	c := srv.Options().Config

	// set a value
	fmt.Println("Value of key.subkey: ", c.Get("key", "subkey").String(""))
}
