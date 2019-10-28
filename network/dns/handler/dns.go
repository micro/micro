package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	dns "github.com/micro/micro/network/dns/proto/dns"
	"github.com/micro/micro/network/dns/provider"
)

// DNS handles incoming gRPC requests
type DNS struct {
	provider provider.Provider
}

// Advertise adds records of network nodes to DNS
func (d *DNS) Advertise(ctx context.Context, req *dns.AdvertiseRequest, rsp *dns.AdvertiseResponse) error {
	log.Debug("Received Advertise Request")
	return d.provider.Advertise(req.Records...)
}

// Remove removes itself from DNS
func (d *DNS) Remove(ctx context.Context, req *dns.RemoveRequest, rsp *dns.RemoveResponse) error {
	return d.provider.Remove(req.Records...)
}

// Resolve looks up matching records and returns any matches
func (d *DNS) Resolve(ctx context.Context, req *dns.ResolveRequest, rsp *dns.ResolveResponse) error {
	providerResponse, err := d.provider.Resolve(req.Name, req.Type)
	if err != nil {
		return err
	}
	rsp.Records = providerResponse
	return nil
}
