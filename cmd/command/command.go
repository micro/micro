package command

import "github.com/micro/cli/v2"

var commands []*cli.Command

// Register CLI commands
func Register(cmds ...*cli.Command) {
	commands = append(commands, cmds...)
}

// List CLI commands
func List() []*cli.Command {
	return commands
}
