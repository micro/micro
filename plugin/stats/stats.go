// Package stats records request stats
package stats

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/micro/v3/internal/wrapper"
	"github.com/micro/micro/v3/plugin"
	s "github.com/micro/micro/v3/service/server"
)

var (
	Plugin = plugin.NewPlugin(
		plugin.WithName("stats"),
		plugin.WithInit(func(ctx *cli.Context) error {
			// wrap the server
			s.DefaultServer.Init(
				server.WrapHandler(wrapper.HandlerStats()),
			)
			return nil
		}),
	)
)

func init() {
	plugin.Register(Plugin)
}
