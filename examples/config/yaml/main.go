package main

import (
	"fmt"

	yaml "github.com/micro/plugins/v5/config/encoder/yaml"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/config/reader"
	"go-micro.dev/v5/config/reader/json"
	"go-micro.dev/v5/config/source/file"
)

func main() {
	// new yaml encoder
	enc := yaml.NewEncoder()

	// new config
	c, _ := config.NewConfig(
		config.WithReader(
			json.NewReader( // json reader for internal config merge
				reader.WithEncoder(enc),
			),
		),
	)

	// load the config from a file source
	if err := c.Load(file.NewSource(
		file.WithPath("./config.yaml"),
	)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("data", c.Map())

	// define our own host type
	type Host struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	var host Host

	// read a database host
	if err := c.Get("hosts", "database").Scan(&host); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(host.Address, host.Port)
}
