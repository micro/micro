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
// Original source: github.com/micro/go-micro/v3/network/options.go

package network

import (
	"github.com/google/uuid"
	"github.com/micro/micro/v3/internal/network/tunnel"
	tmucp "github.com/micro/micro/v3/internal/network/tunnel/mucp"
	"github.com/micro/micro/v3/service/proxy"
	"github.com/micro/micro/v3/service/proxy/mucp"
	"github.com/micro/micro/v3/service/router"
	regRouter "github.com/micro/micro/v3/service/router/registry"
)

type Option func(*Options)

// Options configure network
type Options struct {
	// Id of the node
	Id string
	// Name of the network
	Name string
	// Address to bind to
	Address string
	// Advertise sets the address to advertise
	Advertise string
	// Nodes is a list of nodes to connect to
	Nodes []string
	// Tunnel is network tunnel
	Tunnel tunnel.Tunnel
	// Router is network router
	Router router.Router
	// Proxy is network proxy
	Proxy proxy.Proxy
}

// Id sets the id of the network node
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Name sets the network name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Address sets the network address
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Advertise sets the address to advertise
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// Nodes is a list of nodes to connect to
func Nodes(n ...string) Option {
	return func(o *Options) {
		o.Nodes = n
	}
}

// Tunnel sets the network tunnel
func Tunnel(t tunnel.Tunnel) Option {
	return func(o *Options) {
		o.Tunnel = t
	}
}

// Router sets the network router
func Router(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// Proxy sets the network proxy
func Proxy(p proxy.Proxy) Option {
	return func(o *Options) {
		o.Proxy = p
	}
}

// DefaultOptions returns network default options
func DefaultOptions() Options {
	return Options{
		Id:      uuid.New().String(),
		Name:    "go.micro",
		Address: ":0",
		Tunnel:  tmucp.NewTunnel(),
		Router:  regRouter.NewRouter(),
		Proxy:   mucp.NewProxy(),
	}
}
