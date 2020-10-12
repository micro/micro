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
// Original source: github.com/micro/go-micro/v3/api/resolver/vpath/vpath.go

// Package vpath resolves using http path and recognised versioned urls
package vpath

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/micro/micro/v3/internal/api/resolver"
)

func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &Resolver{opts: resolver.NewOptions(opts...)}
}

type Resolver struct {
	opts resolver.Options
}

var (
	re = regexp.MustCompile("^v[0-9]+$")
)

func (r *Resolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	if req.URL.Path == "/" {
		return nil, errors.New("unknown name")
	}

	options := resolver.NewResolveOptions(opts...)

	parts := strings.Split(req.URL.Path[1:], "/")
	if len(parts) == 1 {
		return &resolver.Endpoint{
			Name:   r.withPrefix(parts...),
			Host:   req.Host,
			Method: req.Method,
			Path:   req.URL.Path,
			Domain: options.Domain,
		}, nil
	}

	// /v1/foo
	if re.MatchString(parts[0]) {
		return &resolver.Endpoint{
			Name:   r.withPrefix(parts[0:2]...),
			Host:   req.Host,
			Method: req.Method,
			Path:   req.URL.Path,
			Domain: options.Domain,
		}, nil
	}

	return &resolver.Endpoint{
		Name:   r.withPrefix(parts[0]),
		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Domain: options.Domain,
	}, nil
}

func (r *Resolver) String() string {
	return "path"
}

// withPrefix transforms "foo" into "go.micro.api.foo"
func (r *Resolver) withPrefix(parts ...string) string {
	p := r.opts.ServicePrefix
	if len(p) > 0 {
		parts = append([]string{p}, parts...)
	}

	return strings.Join(parts, ".")
}
