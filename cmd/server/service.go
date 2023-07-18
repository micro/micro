package server

import (
	ccli "github.com/urfave/cli/v2"
	"micro.dev/v4/cmd"

	// services
	auth "micro.dev/v4/service/auth/server"
	broker "micro.dev/v4/service/broker/server"
	config "micro.dev/v4/service/config/server"
	events "micro.dev/v4/service/events/server"
	network "micro.dev/v4/service/network/server"
	registry "micro.dev/v4/service/registry/server"
	runtime "micro.dev/v4/service/runtime/server"
	store "micro.dev/v4/service/store/server"
)

type srvCommand struct {
	Name    string
	Command ccli.ActionFunc
	Flags   []ccli.Flag
}

var srvCommands = []srvCommand{
	{
		Name:    "auth",
		Command: auth.Run,
	},
	{
		Name:    "broker",
		Command: broker.Run,
	},
	{
		Name:    "config",
		Command: config.Run,
		Flags:   config.Flags,
	},
	{
		Name:    "events",
		Command: events.Run,
	},
	{
		Name:    "network",
		Command: network.Run,
		Flags:   network.Flags,
	},
	{
		Name:    "registry",
		Command: registry.Run,
	},
	{
		Name:    "runtime",
		Command: runtime.Run,
	},
	{
		Name:    "store",
		Command: store.Run,
	},
}

func init() {
	flags := []ccli.Flag{
		&ccli.StringFlag{
			Name:    "name",
			Usage:   "Name of the service",
			EnvVars: []string{"MICRO_SERVICE_NAME"},
			Value:   "service",
		},
		&ccli.StringFlag{
			Name:    "address",
			Usage:   "Address of the service",
			EnvVars: []string{"MICRO_SERVICE_ADDRESS"},
		},
	}

	subcommands := make([]*ccli.Command, len(srvCommands))

	for i, c := range srvCommands {
		// construct the command
		command := &ccli.Command{
			Name:   c.Name,
			Flags:  c.Flags,
			Action: c.Command,
		}

		// set the command
		subcommands[i] = command
	}

	command := &ccli.Command{
		Name:        "service",
		Flags:       flags,
		Subcommands: subcommands,
	}

	cmd.Register(command)
}
