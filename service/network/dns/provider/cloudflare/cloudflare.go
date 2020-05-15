// Package cloudflare is a dns Provider for cloudflare
package cloudflare

import (
	"context"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	miekdns "github.com/miekg/dns"
	"github.com/pkg/errors"

	log "github.com/micro/go-micro/v2/logger"
	dns "github.com/micro/micro/v2/service/network/dns/proto/dns"
	"github.com/micro/micro/v2/service/network/dns/provider"
)

type cfProvider struct {
	api    *cloudflare.API
	zoneID string
}

// New returns a configured cloudflare DNS provider
func New(apiToken, zoneID string) (provider.Provider, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, err
	}

	return &cfProvider{
		api:    api,
		zoneID: zoneID,
	}, nil
}

func (cf *cfProvider) Advertise(records ...*dns.Record) error {
	for _, r := range records {
		_, err := cf.api.CreateDNSRecord(cf.zoneID, cloudflare.DNSRecord{
			Name:     r.GetName(),
			Content:  r.GetValue(),
			Type:     r.GetType(),
			Priority: int(r.GetPriority()),
			TTL:      1,
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func (cf *cfProvider) Remove(records ...*dns.Record) error {
	existing := make(map[string]map[string]cloudflare.DNSRecord)
	existingRecords, err := cf.api.DNSRecords(cf.zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return err
	}
	for _, e := range existingRecords {
		if _, found := existing[e.Name]; !found {
			existing[e.Name] = make(map[string]cloudflare.DNSRecord)
		}
		existing[e.Name][e.Content] = e
	}
	for _, r := range records {
		if _, found := existing[r.Name]; !found {
			return errors.New("Record " + r.Name + " could not be deleted as it doesn't exist")
		}
		toDelete, found := existing[r.Name][r.Value]
		if !found {
			return errors.New("Record " + r.Name + " with address " + r.Value + " could not be deleted as it doesn't exist")
		}
		err := cf.api.DeleteDNSRecord(cf.zoneID, toDelete.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cf *cfProvider) Resolve(name, recordType string) ([]*dns.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dnstype, found := miekdns.StringToType[recordType]
	if !found {
		return nil, errors.New(recordType + " is not a valid record type")
	}
	m := new(miekdns.Msg)
	m.SetQuestion(miekdns.Fqdn(name), dnstype)
	r, err := miekdns.ExchangeContext(ctx, m, "1.0.0.1:53")
	if err != nil {
		return nil, err
	}
	var response []*dns.Record
	for _, answer := range r.Answer {
		h := answer.Header()
		rec := &dns.Record{
			Name: h.Name,
			Type: miekdns.TypeToString[h.Rrtype],
			Ttl:  answer.Header().Ttl,
		}
		if rec.Type != recordType {
			log.Trace("Tried to look up a " + recordType + " record but got a " + rec.Type)
			continue
		}
		switch rec.Type {
		case "A":
			arecord, _ := answer.(*miekdns.A)
			rec.Value = arecord.A.String()
		case "AAAA":
			aaaarecord := answer.(*miekdns.AAAA)
			rec.Value = aaaarecord.AAAA.String()
		case "TXT":
			txtrecord := answer.(*miekdns.TXT)
			rec.Value = strings.Join(txtrecord.Txt, "")
		case "MX":
			mxrecord := answer.(*miekdns.MX)
			rec.Value = mxrecord.Mx
			rec.Priority = uint32(mxrecord.Preference)
		default:
			return nil, errors.New("Can't handle record type " + rec.Type)
		}
		response = append(response, rec)
	}
	return response, nil
}
