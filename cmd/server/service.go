package server

import (
	ccli "github.com/urfave/cli/v2"
	"micro.dev/v4/cmd"

	// services
	api "micro.dev/v4/service/api/server"
	auth "micro.dev/v4/service/auth/server"
	broker "micro.dev/v4/service/broker/server"
	config "micro.dev/v4/service/config/server"
	events "micro.dev/v4/service/events/server"
	network "micro.dev/v4/service/network/server"
	proxy "micro.dev/v4/service/proxy/server"
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
		Name:    "api",
		Command: api.Run,
		Flags:   api.Flags,
	},
	{
		Name:    "auth",
		Command: auth.Run,
		Flags:   auth.Flags,
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
		Name:    "proxy",
		Command: proxy.Run,
		Flags:   proxy.Flags,
	},
	{
		Name:    "registry",
		Command: registry.Run,
	},
	{
		Name:    "runtime",
		Command: runtime.Run,
		Flags:   runtime.Flags,
	},
	{
		Name:    "store",
		Command: store.Run,
	},
}

func init() {
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
		Subcommands: subcommands,
	}

	cmd.Register(command)
}
