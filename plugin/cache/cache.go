// Package cache is for client request caching
package cache

import (
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v3/internal/wrapper"
	"github.com/micro/micro/v3/plugin"
	c "github.com/micro/micro/v3/service/client"
)

var (
	Plugin = plugin.NewPlugin(
		plugin.WithName("cache"),
		plugin.WithInit(func(ctx *cli.Context) error {
			// wrap the client
			c.DefaultClient = wrapper.CacheClient(c.DefaultClient)
			return nil
		}),
	)
)

func init() {
	plugin.Register(Plugin)
}
