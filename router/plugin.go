package router

import (
	"fmt"

	"github.com/micro/micro/v2/plugin"
)

var (
	defaultManager = plugin.NewManager()
)

// Plugins lists the router plugins
func Plugins() []plugin.Plugin {
	return defaultManager.Plugins()
}

// Register registers an router plugin
func Register(pl plugin.Plugin) error {
	if plugin.IsRegistered(pl) {
		return fmt.Errorf("%s registered globally", pl.String())
	}
	return defaultManager.Register(pl)
}
