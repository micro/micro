package cli

import (
	"github.com/micro/cli"
)

func registryCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "list",
			Usage: "List items in registry",
			Subcommands: []cli.Command{
				{
					Name:   "services",
					Usage:  "List services in registry",
					Action: listServices,
				},
			},
		},
		{
			Name:  "register",
			Usage: "Register an item in the registry",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "Register a service with JSON definition",
					Action: registerService,
				},
			},
		},
		{
			Name:  "deregister",
			Usage: "Deregister an item in the registry",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "Deregister a service with JSON definition",
					Action: deregisterService,
				},
			},
		},
		{
			Name:  "get",
			Usage: "Get item from registry",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "Get service from registry",
					Action: getService,
				},
			},
		},
	}
}

func Commands() []cli.Command {
	commands := []cli.Command{
		{
			Name:        "registry",
			Usage:       "Query registry",
			Subcommands: registryCommands(),
		},
		{
			Name:   "query",
			Usage:  "Query a service method using rpc",
			Action: queryService,
		},
		{
			Name:   "stream",
			Usage:  "Query a service method using streaming rpc",
			Action: streamService,
		},
		{
			Name:   "health",
			Usage:  "Query the health of a service",
			Action: queryHealth,
		},
	}

	return append(commands, registryCommands()...)
}
