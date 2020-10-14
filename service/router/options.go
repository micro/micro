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
// Original source: github.com/micro/micro/v3/router/options.go

package router

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/registry/mdns"
)

// Options are router options
type Options struct {
	// Id is router id
	Id string
	// Address is router address
	Address string
	// Gateway is network gateway
	Gateway string
	// Network is network address
	Network string
	// Registry is the local registry
	Registry registry.Registry
	// Context for additional options
	Context context.Context
	// Cache routes
	Cache bool
}

// Id sets Router Id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Address sets router service address
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Gateway sets network gateway
func Gateway(g string) Option {
	return func(o *Options) {
		o.Gateway = g
	}
}

// Network sets router network
func Network(n string) Option {
	return func(o *Options) {
		o.Network = n
	}
}

// Registry sets the local registry
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Cache the routes
func Cache() Option {
	return func(o *Options) {
		o.Cache = true
	}
}

// DefaultOptions returns router default options
func DefaultOptions() Options {
	return Options{
		Id:       uuid.New().String(),
		Network:  DefaultNetwork,
		Registry: mdns.NewRegistry(),
		Context:  context.Background(),
	}
}

type ReadOptions struct {
	Service string
}

type ReadOption func(o *ReadOptions)

// ReadService sets the service to read from the table
func ReadService(s string) ReadOption {
	return func(o *ReadOptions) {
		o.Service = s
	}
}
