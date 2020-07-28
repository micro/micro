// Package broker is the micro broker
package broker

import (
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/micro/v3/service/broker/client"
)

// DefaultBroker implementation
var DefaultBroker broker.Broker = client.NewBroker()
