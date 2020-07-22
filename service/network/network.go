// Package network implements micro network node
package network

import (
	"github.com/micro/go-micro/v2/network"
)

var (
	// DefaultNetwork implementation
	DefaultNetwork network.Network = network.NewNetwork()
)
