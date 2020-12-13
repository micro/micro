package opentelemetry

// Options for opentelemetry:
type Options struct {
	ServiceName          string // The name of this service
	TraceReporterAddress string // The address of a reporting server
}

type Option func(o *Options)

const (
	defaultServiceName          = "service"
	defaultTraceReporterAddress = "localhost:8080"
)

// DefaultOptions returns default options:
func DefaultOptions() Options {
	return Options{
		ServiceName:          defaultServiceName,
		TraceReporterAddress: defaultTraceReporterAddress,
	}
}

// WithTraceReporterAddress configures the address of the trace reporter:
func WithTraceReporterAddress(traceReporterAddress string) Option {
	return func(o *Options) {
		o.TraceReporterAddress = traceReporterAddress
	}
}

// WithServiceName configures the name of this service:
func WithServiceName(serviceName string) Option {
	return func(o *Options) {
		o.ServiceName = serviceName
	}
}
