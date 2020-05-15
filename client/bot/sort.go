package bot

import (
	"github.com/micro/go-micro/v2/agent/command"
)

type sortedCommands struct {
	commands []command.Command
}

func (s sortedCommands) Len() int {
	return len(s.commands)
}

func (s sortedCommands) Less(i, j int) bool {
	return s.commands[i].String() < s.commands[j].String()
}

func (s sortedCommands) Swap(i, j int) {
	s.commands[i], s.commands[j] = s.commands[j], s.commands[i]
}
