// Package provider lets you abstract away any number of DNS providers
package provider

import (
	dns "github.com/micro/micro/v2/network/dns/proto/dns"
)

// Provider is an interface for interacting with a DNS provider
type Provider interface {
	// Advertise creates records in DNS
	Advertise(...*dns.Record) error
	// Remove removes records from DNS
	Remove(...*dns.Record) error
	// Resolve looks up a record in DNS
	Resolve(name, recordType string) ([]*dns.Record, error)
}
