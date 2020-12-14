package wrapper

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/internal/opentelemetry"
	"github.com/micro/micro/v3/service/server"
	"github.com/opentracing/opentracing-go"
)

// OpenTraceHandler wraps a server handler to perform opentracing:
func OpenTraceHandler() server.HandlerWrapper {

	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {

		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			// Concatenate the operation name:
			operationName := fmt.Sprintf(req.Service() + "." + req.Endpoint())

			// Don't trace calls to debug:
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// Start a span from context:
			span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, opentelemetry.DefaultOpenTracer, operationName)
			defer span.Finish()

			// Make the service call:
			if err := h(newCtx, req, rsp); err != nil {
				// Include error info:
				span.SetBaggageItem("error", err.Error())
			}

			return nil
		}
	}
}
