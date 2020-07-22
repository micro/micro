package config

import (
	"fmt"
	"strings"

	"github.com/micro/cli/v2"
)

var (
	// UserCommands defines a set of commands specific to user config
	UserCommands = []*cli.Command{
		{
			// TODO: implement `micro user config` which outputs the config
			Name:        "config",
			Description: "Manage user related config like id, token, namespace, etc",
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Usage:  "Get a value; micro user config get key",
					Action: get,
				},
				{
					Name:   "set",
					Usage:  "Set a key-val; micro user config set key val",
					Action: set,
				},
				{
					Name:   "delete",
					Usage:  "Delete a value; micro user config delete key",
					Action: del,
				},
			},
		},
	}
)

func get(ctx *cli.Context) error {
	args := ctx.Args()
	key := args.Get(0)
	val := args.Get(1)

	val, err := Get(strings.Split(key, ".")...)
	if err != nil {
		return err
	}

	fmt.Println(val)
	return nil
}

func set(ctx *cli.Context) error {
	args := ctx.Args()
	key := args.Get(0)
	val := args.Get(1)

	return Set(val, strings.Split(key, ".")...)
}

func del(ctx *cli.Context) error {
	args := ctx.Args()
	key := args.Get(0)

	if len(key) == 0 {
		return fmt.Errorf("key cannot be blank")
	}

	// TODO: actually delete the key also
	return Set("", strings.Split(key, ".")...)
}
