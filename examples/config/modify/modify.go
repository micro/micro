package main

import (
	"fmt"
	"os"

	"github.com/micro/plugins/v5/config/encoder/toml"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/config/source"
	"go-micro.dev/v5/config/source/file"
)

func main() {
	// new toml encoder
	t := toml.NewEncoder()

	// create a new config
	c, err := config.NewConfig(
		config.WithSource(
			// create a new file source
			file.NewSource(
				// path of file
				file.WithPath("./example.conf"),
				// specify the toml encoder
				source.WithEncoder(t),
			),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// load the config
	if err := c.Load(); err != nil {
		fmt.Println(err)
		return
	}

	// set a value
	c.Set("foo", "bar")

	// now the hacks begin
	vals := c.Map()

	// encode
	v, err := t.Encode(vals)
	if err != nil {
		fmt.Println(err)
		return
	}

	// write the file
	if err := os.WriteFile("./example.conf", v, 0644); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("wrote update to example.conf")
}
