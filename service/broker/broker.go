// Package broker is the micro broker
package broker

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/service"
)

// DefaultBroker implementation
var DefaultBroker broker.Broker = service.NewBroker()
