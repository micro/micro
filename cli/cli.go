// Package cli is a command line interface
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/micro/cli"
)

var (
	prompt = "\nmicro> "

	commands = map[string]*command{
		"quit": &command{"quit", "Exit the CLI", quit},
		"exit": &command{"exit", "Exit the CLI", quit},
		"call": &command{"call", "Call a service", callService},
		"get": &command{"get", "Get service info", getService},
		"stream": &command{"stream", "Stream a call to a service", streamService},
		"health": &command{"health", "Get service health", queryHealth},
		"stats": &command{"stats", "Get service stats", queryStats},
	}
)

type command struct {
	name  string
	usage string
	exec  func(*cli.Context, []string)
}

func runc(c *cli.Context) {
	commands["help"] = &command{"help", "CLI usage", help}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// micro>
		fmt.Fprint(os.Stdout, prompt)

		if !scanner.Scan() {
			return
		}

		// get vals
		args := scanner.Text()
		args = strings.TrimSpace(args)

		// skip no args
		if len(args) == 0 {
			continue
		}

		parts := strings.Split(args, " ")
		if len(parts) == 0 {
			continue
		}

		if cmd, ok := commands[parts[0]]; ok {
			cmd.exec(c, parts[1:])
		} else {
			fmt.Fprint(os.Stdout, "unknown command")
		}
	}
}

func registryCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "list",
			Usage: "List items in registry",
			Subcommands: []cli.Command{
				{
					Name:   "services",
					Usage:  "List services in registry",
					Action: printer(listServices),
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
					Action: printer(registerService),
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
					Action: printer(deregisterService),
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
					Action: printer(getService),
				},
			},
		},
	}
}

func Commands() []cli.Command {
	commands := []cli.Command{
		{
			Name:   "cli",
			Usage:  "Start the interactive cli",
			Action: runc,
		},
		{
			Name:        "registry",
			Usage:       "Query registry",
			Subcommands: registryCommands(),
		},
		{
			Name:   "call",
			Usage:  "Call a service or function",
			Action: printer(callService),
		},
		{
			Name:   "query",
			Usage:  "Deprecated: Use call instead",
			Action: printer(callService),
		},
		{
			Name:   "stream",
			Usage:  "Create a service or function stream",
			Action: printer(streamService),
		},
		{
			Name:   "health",
			Usage:  "Query the health of a service",
			Action: printer(queryHealth),
		},
		{
			Name:   "stats",
			Usage:  "Query the stats of a service",
			Action: printer(queryStats),
		},
	}

	return append(commands, registryCommands()...)
}
