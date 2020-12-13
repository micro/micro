package opentelemetry

import "github.com/opentracing/opentracing-go"

var (
	// DefaultOpenTracer is what anything using an opentracing.Tracer will access:
	DefaultOpenTracer opentracing.Tracer = new(opentracing.NoopTracer)
)
