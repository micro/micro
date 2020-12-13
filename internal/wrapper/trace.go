package wrapper

import (
	"context"
	"strings"

	"github.com/micro/micro/v3/internal/debug/trace"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/debug"
	"github.com/micro/micro/v3/service/server"
)

type traceWrapper struct {
	client.Client
}

func (c *traceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	newCtx, s := debug.DefaultTracer.Start(ctx, req.Service()+"."+req.Endpoint())

	s.Type = trace.SpanTypeRequestOutbound
	err := c.Client.Call(newCtx, req, rsp, opts...)
	if err != nil {
		s.Metadata["error"] = err.Error()
	}

	// finish the trace
	debug.DefaultTracer.Finish(s)

	return err
}

// TraceCall is a call tracing wrapper
func TraceCall(c client.Client) client.Client {
	return &traceWrapper{
		Client: c,
	}
}

// TraceHandler wraps a server handler to perform tracing
func TraceHandler() server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// don't store traces for debug
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// get the span
			newCtx, s := debug.DefaultTracer.Start(ctx, req.Service()+"."+req.Endpoint())
			s.Type = trace.SpanTypeRequestInbound

			err := h(newCtx, req, rsp)
			if err != nil {
				s.Metadata["error"] = err.Error()
			}

			// finish
			debug.DefaultTracer.Finish(s)

			return err
		}
	}
}
