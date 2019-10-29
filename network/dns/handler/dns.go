package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	"github.com/micro/go-micro/metadata"
	dns "github.com/micro/micro/network/dns/proto/dns"
	"github.com/micro/micro/network/dns/provider"
	"github.com/pkg/errors"
)

// DNS handles incoming gRPC requests
type DNS struct {
	provider    provider.Provider
	bearerToken string
}

// Advertise adds records of network nodes to DNS
func (d *DNS) Advertise(ctx context.Context, req *dns.AdvertiseRequest, rsp *dns.AdvertiseResponse) error {
	log.Trace("Received Advertise Request")
	if err := d.validateMetadata(ctx); err != nil {
		return err
	}
	return d.provider.Advertise(req.Records...)
}

// Remove removes itself from DNS
func (d *DNS) Remove(ctx context.Context, req *dns.RemoveRequest, rsp *dns.RemoveResponse) error {
	log.Trace("Received Remove Request")
	if err := d.validateMetadata(ctx); err != nil {
		return err
	}
	return d.provider.Remove(req.Records...)
}

// Resolve looks up matching records and returns any matches
func (d *DNS) Resolve(ctx context.Context, req *dns.ResolveRequest, rsp *dns.ResolveResponse) error {
	log.Trace("Received Resolve Request")
	if err := d.validateMetadata(ctx); err != nil {
		return err
	}
	providerResponse, err := d.provider.Resolve(req.Name, req.Type)
	if err != nil {
		return err
	}
	rsp.Records = providerResponse
	return nil
}

func (d *DNS) validateMetadata(ctx context.Context) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return errors.New("Denied: error getting request metadata")
	}
	token, found := md["Authoriztion"]
	if !found {
		return errors.New("Denied: Authorization metadata not provided")
	}
	if token != d.bearerToken {
		return errors.New("Denied: Authorization metadata is not valid")
	}
	return nil
}
