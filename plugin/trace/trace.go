// Package trace records traces for the service
package log

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/micro/v3/internal/wrapper"
	"github.com/micro/micro/v3/plugin"
	c "github.com/micro/micro/v3/service/client"
	s "github.com/micro/micro/v3/service/server"
)

var (
	Plugin = plugin.NewPlugin(
		plugin.WithName("trace"),
		plugin.WithInit(func(ctx *cli.Context) error {
			// wrap the client
			c.DefaultClient = wrapper.TraceCall(c.DefaultClient)
			c.DefaultClient = wrapper.FromService(c.DefaultClient)

			// wrap the server
			s.DefaultServer.Init(
				server.WrapHandler(wrapper.TraceHandler()),
			)

			return nil
		}),
	)
)

func init() {
	plugin.Register(Plugin)
}
