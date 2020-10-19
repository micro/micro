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
// Original source: github.com/micro/micro/v3/router/mdns/mdns.go

// Package mdns is an mdns router
package mdns

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/micro/micro/v3/internal/mdns"
	"github.com/micro/micro/v3/service/router"
)

// NewRouter returns an initialized dns router
func NewRouter(opts ...router.Option) router.Router {
	options := router.DefaultOptions()
	for _, o := range opts {
		o(&options)
	}
	if len(options.Network) == 0 {
		options.Network = "micro"
	}
	return &mdnsRouter{options}
}

type mdnsRouter struct {
	options router.Options
}

func (m *mdnsRouter) Init(opts ...router.Option) error {
	for _, o := range opts {
		o(&m.options)
	}
	return nil
}

func (m *mdnsRouter) Options() router.Options {
	return m.options
}

func (m *mdnsRouter) Table() router.Table {
	return nil
}

func (m *mdnsRouter) Lookup(service string, opts ...router.LookupOption) ([]router.Route, error) {
	options := router.NewLookup(opts...)

	// check to see if we have the port provided in the service, e.g. go-micro-srv-foo:8000
	srv, port, err := net.SplitHostPort(service)
	if err != nil {
		srv = service
	}

	// query for the host
	entries := make(chan *mdns.ServiceEntry)

	p := mdns.DefaultParams(srv)
	p.Timeout = time.Millisecond * 100
	p.Entries = entries

	// check if we're using our own network
	if len(options.Network) > 0 {
		p.Domain = options.Network
	}

	// do the query
	if err := mdns.Query(p); err != nil {
		return nil, err
	}

	var routes []router.Route

	// compose the routes based on the entries
	for e := range entries {
		addr := e.Host
		// prefer ipv4 addrs
		if len(e.AddrV4) > 0 {
			addr = e.AddrV4.String()
			// else use ipv6
		} else if len(e.AddrV6) > 0 {
			addr = "[" + e.AddrV6.String() + "]"
		} else if len(addr) == 0 {
			continue
		}

		pt := 443

		if e.Port > 0 {
			pt = e.Port
		}

		// set the port
		if len(port) > 0 {
			pt, _ = strconv.Atoi(port)
		}

		routes = append(routes, router.Route{
			Service: service,
			Address: fmt.Sprintf("%s:%d", addr, pt),
			Network: p.Domain,
		})
	}

	return routes, nil
}

func (m *mdnsRouter) Watch(opts ...router.WatchOption) (router.Watcher, error) {
	return nil, nil
}

func (m *mdnsRouter) Close() error {
	return nil
}

func (m *mdnsRouter) String() string {
	return "mdns"
}
