package server

import (
	"fmt"

	"github.com/micro/micro/v3/plugin"
)

var (
	defaultManager = plugin.NewManager()
)

// Plugins lists the server plugins
func Plugins() []plugin.Plugin {
	return defaultManager.Plugins()
}

// Register registers an server plugin
func Register(pl plugin.Plugin) error {
	if plugin.IsRegistered(pl) {
		return fmt.Errorf("%s registered globally", pl.String())
	}
	return defaultManager.Register(pl)
}
