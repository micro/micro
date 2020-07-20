// Package cli is a command line interface
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/micro/cli/v2"
)

var (
	prompt = "micro> "

	commands = map[string]*command{
		"services": {"services", "List services in the registry", listServices},
		"quit":     {"quit", "Exit the CLI", quit},
		"exit":     {"exit", "Exit the CLI", quit},
		"call":     {"call", "Call a service", callService},
		"stream":   {"stream", "Stream a call to a service", streamService},
		"health":   {"health", "Get service health", queryHealth},
		"stats":    {"stats", "Get service stats", queryStats},
	}
)

type command struct {
	name  string
	usage string
	exec  exec
}

func Run(c *cli.Context) error {
	commands["help"] = &command{"help", "CLI usage", help}
	alias := map[string]string{
		"?":  "help",
		"ls": "list",
	}

	r, err := readline.New(prompt)
	if err != nil {
		// TODO return err
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}
	defer r.Close()

	for {
		args, err := r.Readline()
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			return err
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
				// TODO return err
				println(err.Error())
				continue
			}
			println(string(rsp))
		} else {
			// TODO return err
			println("unknown command")
		}
	}
	return nil
}

//Commands for micro calling action
func Commands() []*cli.Command {
	commands := []*cli.Command{
		{
			Name:   "cli",
			Usage:  "Run the interactive CLI",
			Action: Run,
		},
		{
			Name:   "call",
			Usage:  "Call a service e.g micro call greeter Say.Hello '{\"name\": \"John\"}",
			Action: Print(callService),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Usage:   "Set the address of the service instance to call",
					EnvVars: []string{"MICRO_ADDRESS"},
				},
				&cli.StringFlag{
					Name:    "output, o",
					Usage:   "Set the output format; json (default), raw",
					EnvVars: []string{"MICRO_OUTPUT"},
				},
				&cli.StringSliceFlag{
					Name:    "metadata",
					Usage:   "A list of key-value pairs to be forwarded as metadata",
					EnvVars: []string{"MICRO_METADATA"},
				},
			},
		},
		{
			Name:   "stream",
			Usage:  "Create a service stream",
			Action: Print(streamService),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output, o",
					Usage:   "Set the output format; json (default), raw",
					EnvVars: []string{"MICRO_OUTPUT"},
				},
				&cli.StringSliceFlag{
					Name:    "metadata",
					Usage:   "A list of key-value pairs to be forwarded as metadata",
					EnvVars: []string{"MICRO_METADATA"},
				},
			},
		},
		{
			Name:   "stats",
			Usage:  "Query the stats of a service",
			Action: Print(queryStats),
		},
		{
			Name:   "env",
			Usage:  "Get/set micro cli environment",
			Action: Print(listEnvs),
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Usage:  "Get the currently selected environment",
					Action: Print(getEnv),
				},
				{
					Name:   "set",
					Usage:  "Set the environment to use for subsequent commands",
					Action: Print(setEnv),
				},
				{
					Name:   "add",
					Usage:  "Add a new environment `micro env add foo 127.0.0.1:8081`",
					Action: Print(addEnv),
				},
				{
					Name:   "del",
					Usage:  "Delete an environment from your list",
					Action: Print(delEnv),
				},
			},
		},
		{
			Name:   "services",
			Usage:  "List micro services",
			Action: Print(listServices),
		},
	}

	return commands
}
