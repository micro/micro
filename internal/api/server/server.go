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
// Original source: github.com/micro/go-micro/v3/api/server/server.go

package server

import (
	"crypto/tls"
	"net/http"

	"github.com/micro/micro/v3/internal/api/resolver"
	"github.com/micro/micro/v3/internal/api/server/acme"
)

// Server serves api requests
type Server interface {
	Address() string
	Init(opts ...Option) error
	Handle(path string, handler http.Handler)
	Start() error
	Stop() error
}

type Option func(o *Options)

type Options struct {
	EnableACME   bool
	EnableCORS   bool
	ACMEProvider acme.Provider
	EnableTLS    bool
	ACMEHosts    []string
	TLSConfig    *tls.Config
	Resolver     resolver.Resolver
	Wrappers     []Wrapper
}

type Wrapper func(h http.Handler) http.Handler

func WrapHandler(w ...Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w...)
	}
}

func EnableCORS(b bool) Option {
	return func(o *Options) {
		o.EnableCORS = b
	}
}

func EnableACME(b bool) Option {
	return func(o *Options) {
		o.EnableACME = b
	}
}

func ACMEHosts(hosts ...string) Option {
	return func(o *Options) {
		o.ACMEHosts = hosts
	}
}

func ACMEProvider(p acme.Provider) Option {
	return func(o *Options) {
		o.ACMEProvider = p
	}
}

func EnableTLS(b bool) Option {
	return func(o *Options) {
		o.EnableTLS = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func Resolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}
