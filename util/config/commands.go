package config

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	// Commands defines a set of commands for local config
	Commands = []*cli.Command{
		{
			Name:   "get",
			Usage:  "Get a value by specifying [key] as an arg",
			Action: get,
		},
		{
			Name:   "set",
			Usage:  "Set a key-val using [key] [value] as args",
			Action: set,
		},
		{
			Name:   "delete",
			Usage:  "Delete a value using [key] as an arg",
			Action: del,
		},
	}
)

func get(ctx *cli.Context) error {
	args := ctx.Args()
	key := args.Get(0)
	val := args.Get(1)

	val, err := Get(key)
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

	return Set(key, val)
}

func del(ctx *cli.Context) error {
	args := ctx.Args()
	key := args.Get(0)

	if len(key) == 0 {
		return errors.New("key cannot be blank")
	}

	// TODO: actually delete the key also
	return Set(key, "")
}
