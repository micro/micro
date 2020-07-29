package main

import (
	"fmt"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.config-read"),
	)

	// get a value
	c := config.DefaultConfig

	for {
		fmt.Println("Value of key.subkey: ", c.Get("key", "subkey").String(""))
		time.Sleep(time.Second)
	}
}
