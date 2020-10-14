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
// Original source: github.com/micro/micro/v3/router/router.go

package router

import (
	"errors"
	"hash/fnv"
)

var (
	// DefaultRouter implementation
	DefaultRouter Router
)

var (
	// DefaultLink is default network link
	DefaultLink = "local"
	// DefaultLocalMetric is default route cost for a local route
	DefaultMetric int64 = 1
	// DefaultNetwork is default micro network
	DefaultNetwork = "micro"
	// ErrRouteNotFound is returned when no route was found in the routing table
	ErrRouteNotFound = errors.New("route not found")
	// ErrDuplicateRoute is returned when the route already exists
	ErrDuplicateRoute = errors.New("duplicate route")
)

// Router is an interface for a routing control plane
type Router interface {
	// Init initializes the router with options
	Init(...Option) error
	// Options returns the router options
	Options() Options
	// The routing table
	Table() Table
	// Lookup queries routes in the routing table
	Lookup(service string, opts ...LookupOption) ([]Route, error)
	// Watch returns a watcher which tracks updates to the routing table
	Watch(opts ...WatchOption) (Watcher, error)
	// Close the router
	Close() error
	// Returns the router implementation
	String() string
}

// Table is an interface for routing table
type Table interface {
	// Create new route in the routing table
	Create(Route) error
	// Delete existing route from the routing table
	Delete(Route) error
	// Update route in the routing table
	Update(Route) error
	// Read is for querying the table
	Read(...ReadOption) ([]Route, error)
}

// Option used by the router
type Option func(*Options)

// StatusCode defines router status
type StatusCode int

const (
	// Running means the router is up and running
	Running StatusCode = iota
	// Stopped means the router has been stopped
	Stopped
	// Error means the router has encountered error
	Error
)

// Route is a network route
type Route struct {
	// Service is destination service name
	Service string
	// Address is service node address
	Address string
	// Gateway is route gateway
	Gateway string
	// Network is network address
	Network string
	// Router is router id
	Router string
	// Link is network link
	Link string
	// Metric is the route cost metric
	Metric int64
	// Metadata for the route
	Metadata map[string]string
}

// Hash returns route hash sum.
func (r *Route) Hash() uint64 {
	h := fnv.New64()
	h.Reset()
	h.Write([]byte(r.Service + r.Address + r.Gateway + r.Network + r.Router + r.Link))
	return h.Sum64()
}
