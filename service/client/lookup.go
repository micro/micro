// Copyright 2020 Asim Aslam
//
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
// Original source: github.com/micro/go-micro/v3/client/lookup.go

package client

import (
	"context"
	"sort"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/router"
)

// LookupFunc is used to lookup routes for a service
type LookupFunc func(context.Context, Request, CallOptions) ([]string, error)

// LookupRoute for a request using the router and then choose one using the selector
func LookupRoute(ctx context.Context, req Request, opts CallOptions) ([]string, error) {
	// check to see if an address was provided as a call option
	if len(opts.Address) > 0 {
		return opts.Address, nil
	}

	// construct the router query
	query := []router.LookupOption{}

	// if a custom network was requested, pass this to the router. By default the router will use it's
	// own network, which is set during initialisation.
	if len(opts.Network) > 0 {
		query = append(query, router.LookupNetwork(opts.Network))
	}

	// lookup the routes which can be used to execute the request
	routes, err := opts.Router.Lookup(req.Service(), query...)
	if err == router.ErrRouteNotFound {
		return nil, errors.InternalServerError("go.micro.client", "service %s: %s", req.Service(), err.Error())
	} else if err != nil {
		return nil, errors.InternalServerError("go.micro.client", "error getting next %s node: %s", req.Service(), err.Error())
	}

	// sort by lowest metric first
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Metric < routes[j].Metric
	})

	var addrs []string

	for _, route := range routes {
		addrs = append(addrs, route.Address)
	}

	return addrs, nil
}
