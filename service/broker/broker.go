// Package broker is the micro broker
package broker

import (
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/micro/v3/service/broker/client"
)

// DefaultBroker implementation
var DefaultBroker broker.Broker = client.NewBroker()

// Publish a message to a topic
func Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	return DefaultBroker.Publish(topic, m, opts...)
}

// Subscribe to a topic
func Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return DefaultBroker.Subscribe(topic, h, opts...)
}
