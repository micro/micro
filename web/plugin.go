package web

import (
	"fmt"

	"github.com/micro/micro/v2/plugin"
)

var (
	defaultManager = plugin.NewManager()
)

// Plugins lists the web plugins
func Plugins() []plugin.Plugin {
	return defaultManager.Plugins()
}

// Register registers an web plugin
func Register(pl plugin.Plugin) error {
	for _, p := range plugin.Plugins() {
		if p.String() == pl.String() {
			return fmt.Errorf("%s registered globally", pl.String())
		}
	}
	return defaultManager.Register(pl)
}
