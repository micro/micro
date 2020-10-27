# Router

The router is the micro routing control plane. It manages intelligent routing on top of the registry.

## Overview

The router is a layer on top of the micro registry which provides intelligent routing based on metrics, feedback 
and other router information. The registry is a dumb database of sorts. The router enables us to build 
a network topology that morphs and flows based on changing traffic patterns.

## Usage

The router is primarily a package in go-micro that is leveraged by the client/selector and proxy. 
The selector does a route Lookup for a service and gets back optimal routes in an ordered list. 
It can provide feedback to the router about how that route performed which can then be used to 
change the metrics associated with this route. In the event the router receives the same routes 
from other routers with different metrics it can take this into consideration.

The router additionally stores the metadata associated with each service and its nodes and can 
use these for dynamic label based routing. Where the path the request take is based on the label 
associated with the service rather than the address.

## Interface

Rough interface reference

```go
type Router interface {
	// lookup a service, additionally provide label filters
	Lookup(service, ...opts) ([]Route, error)
	// Update the metrics for a route
	Update(Route, Metric) error
	// Advertise routes
	Advertise() (<-chan Route, error)
	// Process adverts
	Process(Route) error
	// Get router status
	Status() (Info, error)
}

type Route struct {
	Service string
	Address string
	Metadata map[string]string
	Metric float64
}
```

