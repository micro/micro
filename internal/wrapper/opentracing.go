package wrapper

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/micro/micro/v3/internal/opentelemetry"
	"github.com/micro/micro/v3/internal/opentelemetry/jaeger"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
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

			// Start a span from context:
			span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, opentelemetry.DefaultOpenTracer, operationName)
			span.SetTag("req.endpoint", req.Endpoint())
			span.SetTag("req.method", req.Method())
			span.SetTag("req.service", req.Service())
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

type wrapper struct {
	handler http.Handler
	tracer  opentracing.Tracer
}

func (w *wrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// We'll use this for the operation name:
	operationName := r.RequestURI

	// Start a span:
	span, newCtx := opentracing.StartSpanFromContextWithTracer(r.Context(), w.tracer, operationName)
	ext.SamplingPriority.Set(span, 1)
	defer span.Finish()

	w.handler.ServeHTTP(rw, r.WithContext(newCtx))
}

// HTTPWrapper returns an HTTP handler wrapper:
func HTTPWrapper(h http.Handler) http.Handler {

	logger.Infof("Preparing an OpenTelemetry HTTPWrapper (%s)", reflect.TypeOf(opentelemetry.DefaultOpenTracer))

	// Shouldn't have to do this (profile isn't giving us the correct tracer):
	openTracer, _, _ := jaeger.New(
		opentelemetry.WithServiceName("API"),
		opentelemetry.WithTraceReporterAddress("localhost:6831"),
	)

	return &wrapper{
		handler: h,
		tracer:  openTracer,
	}
}
