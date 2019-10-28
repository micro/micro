// Package cloudflare is a dns Provider for cloudflare
package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"

	dns "github.com/micro/micro/network/dns/proto/dns"
	"github.com/micro/micro/network/dns/provider"
)

type cfProvider struct {
	api *cloudflare.API
}

// New returns a configured cloudflare DNS provider
func New() (provider.Provider, error) {
	api, _ := cloudflare.NewWithAPIToken("token")

	return &cfProvider{
		api: api,
	}, nil
}

func (cf *cfProvider) Advertise(records ...*dns.Record) error {
	for _, r := range records {
		cf.api.CreateDNSRecord("zoneid", cloudflare.DNSRecord{
			Name:     r.GetName(),
			Content:  r.GetValue(),
			Type:     r.GetType(),
			Priority: int(r.GetPriority()),
		})
	}
	return errors.New("not implemented")
}

func (cf *cfProvider) Remove(records ...*dns.Record) error {
	return errors.New("not implemented")
}

func (cf *cfProvider) Resolve(name, recordType string) ([]*dns.Record, error) {
	return nil, errors.New("not implemented")
}
