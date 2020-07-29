// Package cli is a command line interface
package cli

import (
	"fmt"
	"os"
	osexec "os/exec"
	"strings"

	"github.com/chzyer/readline"
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
)

var (
	prompt = "micro> "

	// TODO: only run fixed set of commands for security purposes
	commands = map[string]*command{}
)

type command struct {
	name  string
	usage string
	exec  util.Exec
}

func Run(c *cli.Context) error {
	// take the first arg as the binary
	binary := os.Args[0]

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

		cmd := osexec.Command(binary, parts...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(string(err.(*osexec.ExitError).Stderr))
		}
	}

	return nil
}

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "cli",
			Usage:  "Run the interactive CLI",
			Action: Run,
		},
		&cli.Command{
			Name:   "call",
			Usage:  "Call a service e.g micro call greeter Say.Hello '{\"name\": \"John\"}",
			Action: util.Print(callService),
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
		&cli.Command{
			Name:   "stream",
			Usage:  "Create a service stream",
			Action: util.Print(streamService),
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
		&cli.Command{
			Name:   "stats",
			Usage:  "Query the stats of a service",
			Action: util.Print(queryStats),
		},
		&cli.Command{
			Name:   "env",
			Usage:  "Get/set micro cli environment",
			Action: util.Print(listEnvs),
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Usage:  "Get the currently selected environment",
					Action: util.Print(getEnv),
				},
				{
					Name:   "set",
					Usage:  "Set the environment to use for subsequent commands",
					Action: util.Print(setEnv),
				},
				{
					Name:   "add",
					Usage:  "Add a new environment `micro env add foo 127.0.0.1:8081`",
					Action: util.Print(addEnv),
				},
				{
					Name:   "del",
					Usage:  "Delete an environment from your list",
					Action: util.Print(delEnv),
				},
			},
		},
		&cli.Command{
			Name:   "services",
			Usage:  "List micro services",
			Action: util.Print(listServices),
		},
	)
}
