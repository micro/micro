package runtime

import (
	"context"
	"io"

	"github.com/micro/go-micro/v3/runtime"
)

type (
	// CreateOptions is an alias for CreateOptions
	CreateOptions = runtime.CreateOptions
	// ReadOptions is an alias for ReadOptions
	ReadOptions = runtime.ReadOptions
	// UpdateOptions is an alias for UpdateOptions
	UpdateOptions = runtime.UpdateOptions
	// DeleteOptions is an alias for DeleteOptions
	DeleteOptions = runtime.DeleteOptions
	// LogsOptions is an alias for LogsOptions
	LogsOptions = runtime.LogsOptions
	// CreateOption is an alias for CreateOption
	CreateOption = runtime.CreateOption
	// ReadOption is an alias for ReadOption
	ReadOption = runtime.ReadOption
	// UpdateOption is an alias for UpdateOption
	UpdateOption = runtime.UpdateOption
	// DeleteOption is an alias for DeleteOption
	DeleteOption = runtime.DeleteOption
	// LogsOption is an alias for LogsOption
	LogsOption = runtime.LogsOption
)

// CreateType sets the type of service to create
func CreateType(t string) CreateOption {
	return func(o *CreateOptions) {
		o.Type = t
	}
}

// CreateImage sets the image to use
func CreateImage(img string) CreateOption {
	return func(o *CreateOptions) {
		o.Image = img
	}
}

// CreateNamespace sets the namespace
func CreateNamespace(ns string) CreateOption {
	return func(o *CreateOptions) {
		o.Namespace = ns
	}
}

// CreateContext sets the context
func CreateContext(ctx context.Context) CreateOption {
	return func(o *CreateOptions) {
		o.Context = ctx
	}
}

// CreateEntrypoint sets the entrypoint
func CreateEntrypoint(e string) CreateOption {
	return func(o *CreateOptions) {
		o.Entrypoint = e
	}
}

// WithSecret sets a secret to provide the service with
func WithSecret(key, value string) CreateOption {
	return func(o *CreateOptions) {
		if o.Secrets == nil {
			o.Secrets = map[string]string{key: value}
		} else {
			o.Secrets[key] = value
		}
	}
}

// WithCommand specifies the command to execute
func WithCommand(cmd ...string) CreateOption {
	return func(o *CreateOptions) {
		// set command
		o.Command = cmd
	}
}

// WithArgs specifies the command to execute
func WithArgs(args ...string) CreateOption {
	return func(o *CreateOptions) {
		// set command
		o.Args = args
	}
}

// WithRetries sets the max retries attemps
func WithRetries(retries int) CreateOption {
	return func(o *CreateOptions) {
		o.Retries = retries
	}
}

// WithEnv sets the created service environment
func WithEnv(env []string) CreateOption {
	return func(o *CreateOptions) {
		o.Env = env
	}
}

// WithVolume adds a volume to be mounted
func WithVolume(name, path string) CreateOption {
	return func(o *CreateOptions) {
		if o.Volumes == nil {
			o.Volumes = map[string]string{name: path}
		} else {
			o.Volumes[name] = path
		}
	}
}

// WithServiceAccount sets the ServiceAccount
func WithServiceAccount(s string) CreateOption {
	return func(o *CreateOptions) {
		o.ServiceAccount = s
	}
}

// WithOutput sets the arg output
func WithOutput(out io.Writer) CreateOption {
	return func(o *CreateOptions) {
		o.Output = out
	}
}

// WithPort sets the port to expose
func WithPort(p string) CreateOption {
	return func(o *CreateOptions) {
		o.Port = p
	}
}

// ResourceLimits sets the resources for the service to use
func ResourceLimits(r *runtime.Resources) CreateOption {
	return func(o *CreateOptions) {
		o.Resources = r
	}
}

// ReadService returns services with the given name
func ReadService(service string) ReadOption {
	return func(o *ReadOptions) {
		o.Service = service
	}
}

// ReadVersion confifgures service version
func ReadVersion(version string) ReadOption {
	return func(o *ReadOptions) {
		o.Version = version
	}
}

// ReadType returns services of the given type
func ReadType(t string) ReadOption {
	return func(o *ReadOptions) {
		o.Type = t
	}
}

// ReadNamespace sets the namespace
func ReadNamespace(ns string) ReadOption {
	return func(o *ReadOptions) {
		o.Namespace = ns
	}
}

// ReadContext sets the context
func ReadContext(ctx context.Context) ReadOption {
	return func(o *ReadOptions) {
		o.Context = ctx
	}
}

// UpdateSecret sets a secret to provide the service with
func UpdateSecret(key, value string) UpdateOption {
	return func(o *UpdateOptions) {
		if o.Secrets == nil {
			o.Secrets = map[string]string{key: value}
		} else {
			o.Secrets[key] = value
		}
	}
}

// UpdateEntrypoint sets the entrypoint
func UpdateEntrypoint(e string) UpdateOption {
	return func(o *UpdateOptions) {
		o.Entrypoint = e
	}
}

// UpdateNamespace sets the namespace
func UpdateNamespace(ns string) UpdateOption {
	return func(o *UpdateOptions) {
		o.Namespace = ns
	}
}

// UpdateContext sets the context
func UpdateContext(ctx context.Context) UpdateOption {
	return func(o *UpdateOptions) {
		o.Context = ctx
	}
}

// DeleteNamespace sets the namespace
func DeleteNamespace(ns string) DeleteOption {
	return func(o *DeleteOptions) {
		o.Namespace = ns
	}
}

// DeleteContext sets the context
func DeleteContext(ctx context.Context) DeleteOption {
	return func(o *DeleteOptions) {
		o.Context = ctx
	}
}

// LogsCount confiures how many existing lines to show
func LogsCount(count int64) LogsOption {
	return func(l *LogsOptions) {
		l.Count = count
	}
}

// LogsStream configures whether to stream new lines
func LogsStream(stream bool) LogsOption {
	return func(l *LogsOptions) {
		l.Stream = stream
	}
}

// LogsNamespace sets the namespace
func LogsNamespace(ns string) LogsOption {
	return func(o *LogsOptions) {
		o.Namespace = ns
	}
}

// LogsContext sets the context
func LogsContext(ctx context.Context) LogsOption {
	return func(o *LogsOptions) {
		o.Context = ctx
	}
}
