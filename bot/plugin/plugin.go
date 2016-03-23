package plugin

import (
	"github.com/micro/cli"
)

var (
	Plugins = map[string]Plugin{}
)

// Plugin is an interface for sources which
// provide a way to communicate with the bot.
// Slack, HipChat, XMPP, etc.
type Plugin interface {
	Flags() []cli.Flag
	Init(*cli.Context) error
	Start() error
	Stop() error
	String() string
}
