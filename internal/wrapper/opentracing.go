package wrapper

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/internal/opentelemetry"
	"github.com/micro/micro/v3/service/server"
)

// OpenTraceHandler wraps a server handler to perform opentracing:
func OpenTraceHandler() server.HandlerWrapper {

	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {

		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			// Don't trace calls to debug:
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// Concatenate the operation name:
			operationName := fmt.Sprintf(req.Service() + "." + req.Endpoint())

			// Get a span:
			span := opentelemetry.DefaultOpenTracer.StartSpan(operationName)

			// Add operation metadata:
			span.SetOperationName(req.Service() + "." + req.Endpoint())

			// Make the service call:
			err := h(ctx, req, rsp)
			if err != nil {
				// Include error info:
				span.SetBaggageItem("error", err.Error())
			}

			// finish
			span.Finish()

			return err
		}
	}
}
