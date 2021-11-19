package wrapper

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/micro/micro/v3/service/client"
	mmd "github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/util/opentelemetry"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
			md, ok := mmd.FromContext(ctx)
			if !ok {
				md = mmd.Metadata{}
			}
			spanCtx, err := opentelemetry.DefaultOpenTracer.Extract(opentracing.TextMap, opentelemetry.MicroMetadataReaderWriter{md})
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				logger.Errorf("Error reconstructing span %s", err)
			}

			// Start a span from context:
			span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, opentelemetry.DefaultOpenTracer, operationName, opentracing.ChildOf(spanCtx), ext.SpanKindRPCServer)
			// TODO remove me
			ext.SamplingPriority.Set(span, 1)
			defer span.Finish()

			// Make the service call, and include error info (if any):
			if err := h(newCtx, req, rsp); err != nil {
				span.SetBaggageItem("error", err.Error())
				return err
			}

			return nil
		}
	}
}

type httpWrapper struct {
	handler http.Handler
}

func (hw *httpWrapper) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	// We'll use this for the operation name:
	operationName := req.RequestURI

	// Initialise a statusRecorder with an assumed 200 status:
	statusRecorder := &statusRecorder{rsp, 200}

	// Start a span:
	span, newCtx := opentracing.StartSpanFromContext(req.Context(), operationName, ext.SpanKindRPCServer)
	ext.SamplingPriority.Set(span, 1)
	defer span.Finish()

	// Handle the request:
	hw.handler.ServeHTTP(statusRecorder, req.WithContext(newCtx))

	// Add trace metadata:
	span.SetTag("req.method", req.Method)
	span.SetTag("rsp.code", statusRecorder.statusCode)

	// Error text:
	if statusRecorder.statusCode >= 500 {
		span.SetBaggageItem("error", http.StatusText(statusRecorder.statusCode))
	}
}

// HTTPWrapper returns an HTTP handler wrapper:
func HTTPWrapper(h http.Handler) http.Handler {
	return &httpWrapper{
		handler: h,
	}
}

// statusRecorder wraps http.ResponseWriter so we can actually capture the status code:
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// Hijack returns the underlying connection
func (sr *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return sr.ResponseWriter.(http.Hijacker).Hijack()
}

// WriteHeader is where we capture the status:
func (sr *statusRecorder) WriteHeader(statusCode int) {
	sr.statusCode = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}

type opentraceWrapper struct {
	client.Client
}

func (o *opentraceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	var span opentracing.Span
	ctx, span = o.wrapContext(ctx, req, opts...)
	defer span.Finish()
	err := o.Client.Call(ctx, req, rsp, opts...)
	if err != nil {
		span.SetBaggageItem("error", err.Error())
		return err
	}
	return nil

}

func (o *opentraceWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	var span opentracing.Span
	ctx, span = o.wrapContext(ctx, req, opts...)
	s, err := o.Client.Stream(ctx, req, opts...)
	if err != nil {
		span.SetBaggageItem("error", err.Error())
		span.Finish()
		return s, err
	}
	go func() {
		<-s.Context().Done()
		span.Finish()
	}()
	return s, nil

}

func (o *opentraceWrapper) wrapContext(ctx context.Context, req client.Request, opts ...client.CallOption) (context.Context, opentracing.Span) {
	// set the open tracing headers
	md := mmd.Metadata{}
	operationName := fmt.Sprintf(req.Service() + "." + req.Endpoint())
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, opentelemetry.DefaultOpenTracer, operationName, ext.SpanKindRPCClient)
	if err := opentelemetry.DefaultOpenTracer.Inject(span.Context(), opentracing.TextMap, opentelemetry.MicroMetadataReaderWriter{md}); err != nil {
		logger.Errorf("Error injecting span %s", err)
	}
	ctx = mmd.MergeContext(newCtx, md, true)

	return ctx, span
}

// OpentraceClient wraps requests with the open tracing headers
func OpentraceClient(c client.Client) client.Client {
	return &opentraceWrapper{c}
}
