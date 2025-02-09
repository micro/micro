package server

import (
	"github.com/micro/micro/v5/cmd"
	ccli "github.com/urfave/cli/v2"

	// services
	auth "github.com/micro/micro/v5/service/auth/server"
	broker "github.com/micro/micro/v5/service/broker/server"
	config "github.com/micro/micro/v5/service/config/server"
	events "github.com/micro/micro/v5/service/events/server"
	network "github.com/micro/micro/v5/service/network/server"
	registry "github.com/micro/micro/v5/service/registry/server"
	runtime "github.com/micro/micro/v5/service/runtime/server"
	store "github.com/micro/micro/v5/service/store/server"
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
