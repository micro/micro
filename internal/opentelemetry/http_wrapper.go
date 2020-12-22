package opentelemetry

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/micro/micro/v3/service/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type wrapper struct {
	handler http.Handler
}

func (w *wrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Concatenate the operation name:
	operationName := fmt.Sprintf("API.%s", r.URL.RawPath)

	// Start a span:
	span, newCtx := opentracing.StartSpanFromContextWithTracer(r.Context(), DefaultOpenTracer, operationName)
	ext.SamplingPriority.Set(span, 1)
	defer span.Finish()

	logger.Warnf("OpenTelemetry Wrapping (%s)", reflect.TypeOf(DefaultOpenTracer))
	w.handler.ServeHTTP(rw, r.WithContext(newCtx))
}

// HTTPWrapper returns an HTTP handler wrapper:
func HTTPWrapper(h http.Handler) http.Handler {
	return &wrapper{
		handler: h,
	}
}
