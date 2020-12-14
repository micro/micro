package wrapper

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/internal/opentelemetry"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	"github.com/opentracing/opentracing-go"
)

// OpenTraceHandler wraps a server handler to perform opentracing:
func OpenTraceHandler() server.HandlerWrapper {

	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {

		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			var span opentracing.Span

			// Concatenate the operation name:
			operationName := fmt.Sprintf(req.Service() + "." + req.Endpoint())

			// Don't trace calls to debug:
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// Extract context and use it to create a child span::
			// spanContext, err := opentelemetry.DefaultOpenTracer.Extract(opentracing.HTTPHeaders, ctx)
			spanContext, err := opentelemetry.ExtractSpanContext(ctx)
			if err != nil {
				logger.Warnf("Unable to extract opentracing context: %v", err)
				span = opentelemetry.DefaultOpenTracer.StartSpan(operationName)
			} else {
				span = opentelemetry.DefaultOpenTracer.StartSpan(operationName, opentracing.ChildOf(spanContext))
			}

			// Inject the context back in:
			newCtx := opentelemetry.InjectSpanContext(ctx, span)

			// Add operation metadata:
			span.SetOperationName(req.Service() + "." + req.Endpoint())

			// Make the service call:
			err = h(newCtx, req, rsp)
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
