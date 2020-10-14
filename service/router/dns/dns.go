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
// Original source: github.com/micro/micro/v3/router/dns/dns.go

package dns

import (
	"fmt"
	"net"
	"strconv"

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
	return &dns{options}
}

type dns struct {
	options router.Options
}

func (d *dns) Init(opts ...router.Option) error {
	for _, o := range opts {
		o(&d.options)
	}
	return nil
}

func (d *dns) Options() router.Options {
	return d.options
}

func (d *dns) Table() router.Table {
	return nil
}

func (d *dns) Close() error {
	return nil
}

func (d *dns) Lookup(service string, opts ...router.LookupOption) ([]router.Route, error) {
	// check to see if we have the port provided in the service, e.g. go-micro-srv-foo:8000
	host, port, err := net.SplitHostPort(service)
	if err == nil {
		// lookup the service using A records
		ips, err := net.LookupHost(host)
		if err != nil {
			return nil, err
		}

		p, _ := strconv.Atoi(port)

		// convert the ip addresses to routes
		result := make([]router.Route, len(ips))
		for i, ip := range ips {
			result[i] = router.Route{
				Service: service,
				Address: fmt.Sprintf("%s:%d", ip, uint16(p)),
			}
		}
		return result, nil
	}

	// we didn't get the port so we'll lookup the service using SRV records. If we can't lookup the
	// service using the SRV record, we return the error.
	_, nodes, err := net.LookupSRV(service, "tcp", d.options.Network)
	if err != nil {
		return nil, err
	}

	// convert the nodes (net services) to routes
	result := make([]router.Route, len(nodes))
	for i, n := range nodes {
		result[i] = router.Route{
			Service: service,
			Address: fmt.Sprintf("%s:%d", n.Target, n.Port),
			Network: d.options.Network,
		}
	}
	return result, nil
}

func (d *dns) Watch(opts ...router.WatchOption) (router.Watcher, error) {
	return nil, nil
}

func (d *dns) String() string {
	return "dns"
}
