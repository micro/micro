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
// Original source: github.com/micro/go-micro/v3/network/resolver/dnssrv/dnssrv.go

// Package dns srv resolves names to dns srv records
package dnssrv

import (
	"fmt"
	"net"

	"github.com/micro/micro/v3/internal/network/resolver"
)

// Resolver is a DNS network resolve
type Resolver struct{}

// Resolve assumes ID is a domain name e.g micro.mu
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	_, addrs, err := net.LookupSRV("network", "udp", name)
	if err != nil {
		return nil, err
	}
	records := make([]*resolver.Record, 0, len(addrs))
	for _, addr := range addrs {
		address := addr.Target
		if addr.Port > 0 {
			address = fmt.Sprintf("%s:%d", addr.Target, addr.Port)
		}
		records = append(records, &resolver.Record{
			Address: address,
		})
	}
	return records, nil
}
