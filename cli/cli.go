// Package cli is a command line interface
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/micro/cli"
)

var (
	prompt = "micro> "

	commands = map[string]*command{
		"quit":       &command{"quit", "Exit the CLI", quit},
		"exit":       &command{"exit", "Exit the CLI", quit},
		"call":       &command{"call", "Call a service", callService},
		"list":       &command{"list", "List services", listServices},
		"get":        &command{"get", "Get service info", getService},
		"stream":     &command{"stream", "Stream a call to a service", streamService},
		"publish":    &command{"publish", "Publish a message to a topic", publish},
		"health":     &command{"health", "Get service health", queryHealth},
		"stats":      &command{"stats", "Get service stats", queryStats},
		"register":   &command{"register", "Register a service", registerService},
		"deregister": &command{"deregister", "Deregister a service", deregisterService},
	}
)

type command struct {
	name  string
	usage string
	exec  exec
}

func runc(c *cli.Context) {
	commands["help"] = &command{"help", "CLI usage", help}
	alias := map[string]string{
		"?":  "help",
		"ls": "list",
	}

	r, err := readline.New(prompt)
	if err != nil {
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}
	defer r.Close()

	for {
		args, err := r.Readline()
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			return
		}

		args = strings.TrimSpace(args)

		// skip no args
		if len(args) == 0 {
			continue
		}

		parts := strings.Split(args, " ")
		if len(parts) == 0 {
			continue
		}

		name := parts[0]

		// get alias
		if n, ok := alias[name]; ok {
			name = n
		}

		if cmd, ok := commands[name]; ok {
			rsp, err := cmd.exec(c, parts[1:])
			if err != nil {
				println(err.Error())
				continue
			}
			println(string(rsp))
		} else {
			println("unknown command")
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
			Name:  "query",
			Usage: "Deprecated: Use call instead",
			Action: func(c *cli.Context) {
				fmt.Println("Deprecated. Use call instead")
				printer(callService)(c)
			},
		},
		{
			Name:   "stream",
			Usage:  "Create a service or function stream",
			Action: printer(streamService),
		},
		{
			Name:   "publish",
			Usage:  "Publish a message to a topic",
			Action: printer(publish),
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
