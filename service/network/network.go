// Package network implements micro network node
package network

import (
	"github.com/micro/go-micro/v2/network"
	"github.com/micro/micro/v2/service/network/client"
)

var (
	// DefaultNetwork implementation
	DefaultNetwork network.Network = client.NewNetwork()
)
