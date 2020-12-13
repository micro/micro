package wrapper

import (
	"context"

	"github.com/micro/micro/v3/service/debug"
	"github.com/micro/micro/v3/service/server"
)

// HandlerStats wraps a server handler to generate request/error stats
func HandlerStats() server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// execute the handler
			err := h(ctx, req, rsp)
			// record the stats
			debug.DefaultStats.Record(err)
			// return the error
			return err
		}
	}
}
