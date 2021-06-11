package opentelemetry

// Options for opentelemetry:
type Options struct {
	SamplingRate         float64 // Percentage of requests to sample (0 = rely on propagated decision)
	ServiceName          string  // The name of this service
	TraceReporterAddress string  // The address of a reporting server
}

type Option func(o *Options)

const (
	defaultSamplingRate         = float64(0)
	defaultServiceName          = "Micro"
	defaultTraceReporterAddress = "localhost:6831"
)

// DefaultOptions returns default options:
func DefaultOptions() Options {
	return Options{
		SamplingRate:         defaultSamplingRate,
		ServiceName:          defaultServiceName,
		TraceReporterAddress: defaultTraceReporterAddress,
	}
}

// WithSamplingRate configures the sampling rate:
func WithSamplingRate(samplingRate float64) Option {
	return func(o *Options) {
		o.SamplingRate = samplingRate
	}
}

// WithServiceName configures the name of this service:
func WithServiceName(serviceName string) Option {
	return func(o *Options) {
		o.ServiceName = serviceName
	}
}

// WithTraceReporterAddress configures the address of the trace reporter:
func WithTraceReporterAddress(traceReporterAddress string) Option {
	return func(o *Options) {
		o.TraceReporterAddress = traceReporterAddress
	}
}
