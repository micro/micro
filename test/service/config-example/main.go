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
			time.Sleep(time.Second)
			fmt.Println("Value of key.subkey: ", config.Get("key", "subkey").String(""))
		}
	}()

	// run the service
	service.Run()
}
