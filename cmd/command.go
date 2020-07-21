package cmd

import "github.com/micro/cli/v2"

var cliCommands []*cli.Command

// Register CLI commands
func Register(cmds ...*cli.Command) {
	cliCommands = append(cliCommands, cmds...)
}
