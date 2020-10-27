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
// Original source: github.com/micro/go-micro/v3/api/api.go

package api

import (
	"errors"
	"regexp"
	"strings"

	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/server"
)

type API interface {
	// Initialise options
	Init(...Option) error
	// Get the options
	Options() Options
	// Register a http handler
	Register(*Endpoint) error
	// Register a route
	Deregister(*Endpoint) error
	// Implementation of api
	String() string
}

type Options struct{}

type Option func(*Options) error

// Endpoint is a mapping between an RPC method and HTTP endpoint
type Endpoint struct {
	// RPC Method e.g. Greeter.Hello
	Name string
	// Description e.g what's this endpoint for
	Description string
	// API Handler e.g rpc, proxy
	Handler string
	// HTTP Host e.g example.com
	Host []string
	// HTTP Methods e.g GET, POST
	Method []string
	// HTTP Path e.g /greeter. Expect POSIX regex
	Path []string
	// Body destination
	// "*" or "" - top level message value
	// "string" - inner message value
	Body string
	// Stream flag
	Stream bool
}

// Service represents an API service
type Service struct {
	// Name of service
	Name string
	// The endpoint for this service
	Endpoint *Endpoint
	// Versions of this service
	Services []*registry.Service
}

// Encode encodes an endpoint to endpoint metadata
func Encode(e *Endpoint) map[string]string {
	if e == nil {
		return nil
	}

	// endpoint map
	ep := make(map[string]string)

	// set vals only if they exist
	set := func(k, v string) {
		if len(v) == 0 {
			return
		}
		ep[k] = v
	}

	set("endpoint", e.Name)
	set("description", e.Description)
	set("handler", e.Handler)
	set("method", strings.Join(e.Method, ","))
	set("path", strings.Join(e.Path, ","))
	set("host", strings.Join(e.Host, ","))

	return ep
}

// Decode decodes endpoint metadata into an endpoint
func Decode(e map[string]string) *Endpoint {
	if e == nil {
		return nil
	}

	return &Endpoint{
		Name:        e["endpoint"],
		Description: e["description"],
		Method:      slice(e["method"]),
		Path:        slice(e["path"]),
		Host:        slice(e["host"]),
		Handler:     e["handler"],
	}
}

// Validate validates an endpoint to guarantee it won't blow up when being served
func Validate(e *Endpoint) error {
	if e == nil {
		return errors.New("endpoint is nil")
	}

	if len(e.Name) == 0 {
		return errors.New("name required")
	}

	for _, p := range e.Path {
		ps := p[0]
		pe := p[len(p)-1]

		if ps == '^' && pe == '$' {
			_, err := regexp.CompilePOSIX(p)
			if err != nil {
				return err
			}
		} else if ps == '^' && pe != '$' {
			return errors.New("invalid path")
		} else if ps != '^' && pe == '$' {
			return errors.New("invalid path")
		}
	}

	if len(e.Handler) == 0 {
		return errors.New("invalid handler")
	}

	return nil
}

func WithEndpoint(e *Endpoint) server.HandlerOption {
	return server.EndpointMetadata(e.Name, Encode(e))
}

func slice(s string) []string {
	var sl []string

	for _, p := range strings.Split(s, ",") {
		if str := strings.TrimSpace(p); len(str) > 0 {
			sl = append(sl, strings.TrimSpace(p))
		}
	}

	return sl
}
