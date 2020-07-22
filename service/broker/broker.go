// Package broker is the micro broker
package broker

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/micro/v2/service/broker/client"
)

// DefaultBroker implementation
var DefaultBroker broker.Broker = client.NewBroker()
