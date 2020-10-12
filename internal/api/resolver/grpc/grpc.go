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
// Original source: github.com/micro/go-micro/v3/api/resolver/grpc/grpc.go

// Package grpc resolves a grpc service like /greeter.Say/Hello to greeter service
package grpc

import (
	"errors"
	"net/http"
	"strings"

	"github.com/micro/micro/v3/internal/api/resolver"
)

type Resolver struct {
	opts resolver.Options
}

func (r *Resolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	// parse options
	options := resolver.NewResolveOptions(opts...)

	// /foo.Bar/Service
	if req.URL.Path == "/" {
		return nil, errors.New("unknown name")
	}
	// [foo.Bar, Service]
	parts := strings.Split(req.URL.Path[1:], "/")
	// [foo, Bar]
	name := strings.Split(parts[0], ".")
	// foo
	return &resolver.Endpoint{
		Name:   strings.Join(name[:len(name)-1], "."),
		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Domain: options.Domain,
	}, nil
}

func (r *Resolver) String() string {
	return "grpc"
}

func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &Resolver{opts: resolver.NewOptions(opts...)}
}
