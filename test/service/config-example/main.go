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
			val, err := config.Get("key.subkey")
			fmt.Println("Value of key.subkey: ", val.String(""), err)
		}
	}()

	// run the service
	service.Run()
}
