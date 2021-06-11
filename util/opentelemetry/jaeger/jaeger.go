package jaeger

import (
	"io"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/util/opentelemetry"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	DefaultReporterAddress = "localhost:6831"
)

// New returns a configured Jaeger tracer:
func New(opts ...opentelemetry.Option) (opentracing.Tracer, io.Closer, error) {
	options := opentelemetry.DefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	logger.Debug("Creating a new Jaeger tracer")

	// Prepare a Jaeger config using our options:
	jaegerConfig := config.Configuration{
		ServiceName: options.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const", // No adaptive sampling or external lookups
			Param: options.SamplingRate,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: options.TraceReporterAddress,
		},
	}

	// Prepare a new Jaeger tracer from this config:
	tracer, closer, err := jaegerConfig.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, nil
}
