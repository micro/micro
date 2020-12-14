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
			var span opentracing.Span

			// Concatenate the operation name:
			operationName := fmt.Sprintf(req.Service() + "." + req.Endpoint())

			// Don't trace calls to debug:
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// Start a span from context:
			span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, opentelemetry.DefaultOpenTracer, operationName)

			// Make the service call:
			if err := h(newCtx, req, rsp); err != nil {
				// Include error info:
				span.SetBaggageItem("error", err.Error())
			}

			// finish
			span.Finish()

			return nil
		}
	}
}

// type openTraceWrapper struct {
// 	client.Client
// }

// func (c *openTraceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
// 	// Concatenate the operation name:
// 	operationName := fmt.Sprintf(req.Service() + "." + req.Endpoint())

// 	span := opentelemetry.DefaultOpenTracer.StartSpan(operationName)
// 	span.SetTag("SAMPLING_PRIORITY", 1)

// 	// Inject the context back in:
// 	newCtx := opentelemetry.InjectSpanContext(ctx, span)

// 	err := c.Client.Call(newCtx, req, rsp, opts...)
// 	if err != nil {
// 		span.SetBaggageItem("error", err.Error())
// 	}

// 	// finish
// 	span.Finish()

// 	return err
// }

// // OpenTraceCall is a call tracing wrapper
// func OpenTraceCall(c client.Client) client.Client {
// 	return &traceWrapper{
// 		Client: c,
// 	}
// }
