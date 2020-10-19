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
// Original source: github.com/micro/go-micro/v3/api/router/router.go

// Package router provides api service routing
package router

import (
	"net/http"

	"github.com/micro/micro/v3/service/api"
)

// Router is used to determine an endpoint for a request
type Router interface {
	// Returns options
	Options() Options
	// Stop the router
	Close() error
	// Endpoint returns an api.Service endpoint or an error if it does not exist
	Endpoint(r *http.Request) (*api.Service, error)
	// Register endpoint in router
	Register(ep *api.Endpoint) error
	// Deregister endpoint from router
	Deregister(ep *api.Endpoint) error
	// Route returns an api.Service route
	Route(r *http.Request) (*api.Service, error)
}
