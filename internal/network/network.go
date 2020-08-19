// Package network provides internal namespaced networking
package network

import (
	"context"

	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/metadata"
)

// Lookup provides a lookup function that checks for namespace as the Micro-Namespace header
func Lookup(ctx context.Context, req client.Request, opts client.CallOptions) ([]string, error) {
	// only set if the value is already nil
	if len(opts.Network) == 0 {
		val, ok := metadata.Get(ctx, "Micro-Namespace")
		if ok {
			// use namespace instead
			opts.Network = val
		}
	}

	// use the standard Lookup function
	return client.LookupRoute(ctx, req, opts)
}
