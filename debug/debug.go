// Package debug implements metrics/logging/introspection/... of go-micro services
package debug

import (
	"fmt"
	"os"

	"github.com/micro/cli"
	"github.com/micro/micro/debug/web"
	"github.com/micro/go-micro"
)

// Commands is called by micro/cmd to populate the micro debug command and subcommands
func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "debug",
		Usage: "Run debug commands",
		Action: func(c *cli.Context) {
			fmt.Fprintf(os.Stderr, "Usage: %s debug web\n", os.Args[0])
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name: "web",
				Usage: "Start the debug web dashboard",
				Action: func(c *cli.Context) {
					web.Run(c)
				},
			},
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}
	return []cli.Command{command}
}
