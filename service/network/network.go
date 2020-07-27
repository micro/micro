// Package network implements micro network node
package network

import (
	"github.com/micro/go-micro/v3/network"
	"github.com/micro/go-micro/v3/network/mucp"
)

var (
	// DefaultNetwork implementation
	DefaultNetwork network.Network = mucp.NewNetwork()
)
