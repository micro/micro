// Package network provides internal namespaced networking
package network

import (
	"context"

	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
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
