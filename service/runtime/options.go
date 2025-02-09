// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/runtime/options.go

package runtime

import (
	"context"
	"io"

	"github.com/micro/micro/v5/service/client"
)

type Option func(o *Options)

// Options configure runtime
type Options struct {
	// Service type to manage
	Type string
	// Client to use when making requests
	Client client.Client
	// Base image to use
	Image string
	// Source of the services repository
	Source string
	// Context to store additional options
	Context context.Context
}

// WithSource sets the base image / repository
func WithSource(src string) Option {
	return func(o *Options) {
		o.Source = src
	}
}

// WithType sets the service type to manage
func WithType(t string) Option {
	return func(o *Options) {
		o.Type = t
	}
}

// WithImage sets the image to use
func WithImage(t string) Option {
	return func(o *Options) {
		o.Image = t
	}
}

// WithClient sets the client to use
func WithClient(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

type CreateOption func(o *CreateOptions)

type ReadOption func(o *ReadOptions)

// CreateOptions configure runtime services
type CreateOptions struct {
	// Command to execut
	Command []string
	// Args to pass into command
	Args []string
	// Environment to configure
	Env []string
	// Entrypoint within the folder (e.g. in the case of a mono-repo)
	Entrypoint string
	// Log output
	Output io.Writer
	// Type of service to create
	Type string
	// Retries before failing deploy
	Retries int
	// Specify the image to use
	Image string
	// Port to expose
	Port string
	// Namespace to create the service in
	Namespace string
	// Specify the context to use
	Context context.Context
	// Secrets to use
	Secrets map[string]string
	// Resources to allocate the service
	Resources *Resources
	// Volumes to mount
	Volumes map[string]string
	// ServiceAccount to start the container with
	ServiceAccount string
	// Number of instances to run
	Instances int
	// Force the service ignore the service status
	Force bool
}

// ReadOptions queries runtime services
type ReadOptions struct {
	// Service name
	Service string
	// Version queries services with given version
	Version string
	// Type of service
	Type string
	// Namespace the service is running in
	Namespace string
	// Specify the context to use
	Context context.Context
}

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

// WithServiceAccount sets the ServiceAccount
func WithServiceAccount(s string) CreateOption {
	return func(o *CreateOptions) {
		o.ServiceAccount = s
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

// WithOutput sets the arg output
func WithOutput(out io.Writer) CreateOption {
	return func(o *CreateOptions) {
		o.Output = out
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

// WithPort sets the port to expose
func WithPort(p string) CreateOption {
	return func(o *CreateOptions) {
		o.Port = p
	}
}

// CreateInstances sets the number of instances
func CreateInstances(v int) CreateOption {
	return func(o *CreateOptions) {
		o.Instances = v
	}
}

// ResourceLimits sets the resources for the service to use
func ResourceLimits(r *Resources) CreateOption {
	return func(o *CreateOptions) {
		o.Resources = r
	}
}

// WithForce sets the sign to force restart the service
func WithForce(f bool) CreateOption {
	return func(o *CreateOptions) {
		o.Force = f
	}
}

// ReadService returns services with the given name
func ReadService(service string) ReadOption {
	return func(o *ReadOptions) {
		o.Service = service
	}
}

// ReadVersion configures service version
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

type UpdateOption func(o *UpdateOptions)

type UpdateOptions struct {
	// Entrypoint within the folder (e.g. in the case of a mono-repo)
	Entrypoint string
	// Namespace the service is running in
	Namespace string
	// Specify the context to use
	Context context.Context
	// Secrets to use
	Secrets map[string]string
	// Number of instances
	Instances int
}

// WithSecret sets a secret to provide the service with
func UpdateSecret(key, value string) UpdateOption {
	return func(o *UpdateOptions) {
		if o.Secrets == nil {
			o.Secrets = map[string]string{key: value}
		} else {
			o.Secrets[key] = value
		}
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

// UpdateEntrypoint sets the entrypoint
func UpdateEntrypoint(e string) UpdateOption {
	return func(o *UpdateOptions) {
		o.Entrypoint = e
	}
}

// UpdateInstances sets the number of instances
func UpdateInstances(v int) UpdateOption {
	return func(o *UpdateOptions) {
		o.Instances = v
	}
}

type DeleteOption func(o *DeleteOptions)

type DeleteOptions struct {
	// Namespace the service is running in
	Namespace string
	// Specify the context to use
	Context context.Context
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

// LogsOption configures runtime logging
type LogsOption func(o *LogsOptions)

// LogsOptions configure runtime logging
type LogsOptions struct {
	// How many existing lines to show
	Count int64
	// Stream new lines?
	Stream bool
	// Namespace the service is running in
	Namespace string
	// Specify the context to use
	Context context.Context
}

// LogsCount configures how many existing lines to show
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
