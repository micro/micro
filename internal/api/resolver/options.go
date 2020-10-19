// Copyright 2020 Asim Aslam
//
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
// Original source: github.com/micro/go-micro/v3/api/resolver/options.go

package resolver

import (
	"github.com/micro/micro/v3/service/registry"
)

type Options struct {
	Handler       string
	ServicePrefix string
}

type Option func(o *Options)

// WithHandler sets the handler being used
func WithHandler(h string) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

// WithServicePrefix sets the ServicePrefix option
func WithServicePrefix(p string) Option {
	return func(o *Options) {
		o.ServicePrefix = p
	}
}

// NewOptions returns new initialised options
func NewOptions(opts ...Option) Options {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ResolveOptions are used when resolving a request
type ResolveOptions struct {
	Domain string
}

// ResolveOption sets an option
type ResolveOption func(*ResolveOptions)

// Domain sets the resolve Domain option
func Domain(n string) ResolveOption {
	return func(o *ResolveOptions) {
		o.Domain = n
	}
}

// NewResolveOptions returns new initialised resolve options
func NewResolveOptions(opts ...ResolveOption) ResolveOptions {
	var options ResolveOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Domain) == 0 {
		options.Domain = registry.DefaultDomain
	}

	return options
}
