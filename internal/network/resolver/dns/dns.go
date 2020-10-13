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
// Original source: github.com/micro/go-micro/v3/network/resolver/dns/dns.go

// Package dns resolves names to dns records
package dns

import (
	"context"
	"net"

	"github.com/micro/micro/v3/internal/network/resolver"
	"github.com/miekg/dns"
)

// Resolver is a DNS network resolve
type Resolver struct {
	// The resolver address to use
	Address string
}

// Resolve assumes ID is a domain name e.g micro.mu
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	host, port, err := net.SplitHostPort(name)
	if err != nil {
		host = name
		port = "8085"
	}

	if len(host) == 0 {
		host = "localhost"
	}

	if len(r.Address) == 0 {
		r.Address = "1.0.0.1:53"
	}

	//nolint:prealloc
	var records []*resolver.Record

	// parsed an actual ip
	if v := net.ParseIP(host); v != nil {
		records = append(records, &resolver.Record{
			Address: net.JoinHostPort(host, port),
		})
		return records, nil
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)
	rec, err := dns.ExchangeContext(context.Background(), m, r.Address)
	if err != nil {
		return nil, err
	}

	for _, answer := range rec.Answer {
		h := answer.Header()
		// check record type matches
		if h.Rrtype != dns.TypeA {
			continue
		}

		arec, _ := answer.(*dns.A)
		addr := arec.A.String()

		// join resolved record with port
		address := net.JoinHostPort(addr, port)
		// append to record set
		records = append(records, &resolver.Record{
			Address: address,
		})
	}

	// no records returned so just best effort it
	if len(records) == 0 {
		records = append(records, &resolver.Record{
			Address: net.JoinHostPort(host, port),
		})
	}

	return records, nil
}
