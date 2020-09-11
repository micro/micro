// Package broker is the micro broker
package broker

import (
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/micro/v3/service/broker/client"
)

// DefaultBroker implementation
var DefaultBroker broker.Broker = client.NewBroker()

type (
	// Handler is an alias for broker.Handler
	Handler = broker.Handler
	// Message is an alias for broker.Message
	Message = broker.Message
	// Subscriber is an alias for broker.Subscriber
	Subscriber = broker.Subscriber
)

// Publish a message to a topic
func Publish(topic string, m *Message, opts ...broker.PublishOption) error {
	return DefaultBroker.Publish(topic, m, opts...)
}

// Subscribe to a topic
func Subscribe(topic string, h Handler, opts ...broker.SubscribeOption) (Subscriber, error) {
	return DefaultBroker.Subscribe(topic, h, opts...)
}
