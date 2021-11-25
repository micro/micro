package opentelemetry

import (
	"strings"

	mmd "github.com/micro/micro/v3/service/context/metadata"
	"github.com/opentracing/opentracing-go"
)

var (
	// DefaultOpenTracer is what anything using an opentracing.Tracer will access:
	DefaultOpenTracer opentracing.Tracer = new(opentracing.NoopTracer)
)

// MicroMetadataReaderWriter satisfies both the opentracing.TextMapReader and
// opentracing.TextMapWriter interfaces.
type MicroMetadataReaderWriter struct {
	mmd.Metadata
}

func (w MicroMetadataReaderWriter) Set(key, val string) {
	// The GRPC HPACK implementation rejects any uppercase keys here.
	//
	// As such, since the HTTP_HEADERS format is case-insensitive anyway, we
	// blindly lowercase the key (which is guaranteed to work in the
	// Inject/Extract sense per the OpenTracing spec).
	key = strings.ToLower(key)
	w.Metadata.Set(key, val)
}

func (w MicroMetadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, v := range w.Metadata {
		if err := handler(k, v); err != nil {
			return err
		}
	}

	return nil
}
