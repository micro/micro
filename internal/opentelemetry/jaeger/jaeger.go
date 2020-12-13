package jaeger

import (
	"github.com/micro/micro/v3/internal/opentelemetry"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// New returns a configured Jaeger tracer:
func New(opts ...opentelemetry.Option) (opentracing.Tracer, error) {
	options := opentelemetry.DefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	// Prepare a Jaeger config using our options:
	jaegerConfig := config.Configuration{
		ServiceName: options.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const", // No adaptive sampling or external lookups
			Param: 0,       // Never randomly decide to trace (only trace if the sampling decision has been propagated to do so)
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: options.TraceReporterAddress,
		},
	}

	// Prepare a new Jaeger tracer from this config:
	tracer, _, err := jaegerConfig.NewTracer()
	if err != nil {
		return nil, err
	}

	return tracer, nil
}
