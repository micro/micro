package server

import (
	"github.com/micro/micro/v3/cmd"
	ccli "github.com/urfave/cli/v2"

	// services
	api "github.com/micro/micro/v3/service/api/server"
	auth "github.com/micro/micro/v3/service/auth/server"
	broker "github.com/micro/micro/v3/service/broker/server"
	config "github.com/micro/micro/v3/service/config/server"
	events "github.com/micro/micro/v3/service/events/server"
	network "github.com/micro/micro/v3/service/network/server"
	proxy "github.com/micro/micro/v3/service/proxy/server"
	registry "github.com/micro/micro/v3/service/registry/server"
	runtime "github.com/micro/micro/v3/service/runtime/server"
	store "github.com/micro/micro/v3/service/store/server"
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
