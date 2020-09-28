package main

import (
	"fmt"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

type keyConfig struct {
	Subkey  string `json:"subkey"`
	Subkey1 int    `json:"subkey1"`
	Subkey2 string `json:"subkey2"`
}

type conf struct {
	Key keyConfig `json:"key"`
}

func main() {
	// get a value
	go func() {
		for {
			time.Sleep(time.Second)
			val, err := config.Get("key.subkey")
			fmt.Println("Value of key.subkey: ", val.String(""), err)

			val, _ = config.Get("key", config.Secret(true))
			c := conf{}
			err = val.Scan(&c.Key)
			fmt.Println("Value of key.subkey1: ", c.Key.Subkey1, err)
			fmt.Println("Value of key.subkey2: ", c.Key.Subkey2)

			val, _ = config.Get("key.subkey3")
			fmt.Println("Value of key.subkey3: ", val.String(""))

			// Test defaults
			val, _ = config.Get("key.subkey_does_not_exist")
			fmt.Println("Default", val.String("Hello"))
		}
	}()

	// run the service
	service.Run()
}
