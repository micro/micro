// Package cli is a command line interface
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
			Name:   "call",
			Usage:  "Call a service or function",
			Action: callService,
		},
		{
			Name:   "query",
			Usage:  "Deprecated: Use call instead",
			Action: callService,
		},
		{
			Name:   "stream",
			Usage:  "Create a service or function stream",
			Action: streamService,
		},
		{
			Name:   "health",
			Usage:  "Query the health of a service",
			Action: queryHealth,
		},
		{
			Name:   "stats",
			Usage:  "Query the stats of a service",
			Action: queryStats,
		},
	}

	return append(commands, registryCommands()...)
}
