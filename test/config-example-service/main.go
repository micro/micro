package main

import (
	"fmt"

	"github.com/micro/go-micro/v2/config"
)

func main() {
	// create a new config
	c, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// set a value
	fmt.Println("Value of key.subkey: ", c.Get("key", "subkey").String(""))
}
