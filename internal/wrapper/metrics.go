package wrapper

import (
	"context"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/metrics"
	"github.com/micro/micro/v3/service/server"
)

// MetricsHandler wraps a server handler to instrument calls
func MetricsHandler() server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			// Don't instrument debug calls:
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// Build some tags to describe the call:
			tags := metrics.Tags{
				"method": req.Method(),
			}

			// Start the clock:
			callTime := time.Now()

			// Run the handlerFunction:
			err := h(ctx, req, rsp)

			// Add a result tag:
			if err != nil {
				tags["result"] = "failure"
			} else {
				tags["result"] = "success"
			}

			// Instrument the result (if the DefaultClient has been configured):
			metrics.Timing("service.handler", time.Since(callTime), tags)

			return err
		}
	}
}
