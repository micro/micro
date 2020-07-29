package main

import (
	"fmt"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

func main() {
	// get a value
	go func() {
		for {
			fmt.Println("Value of key.subkey: ", config.Get("key", "subkey").String(""))
			time.Sleep(time.Second)
		}
	}()

	// run the service
	service.Run()
}
