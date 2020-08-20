package events

import (
	"github.com/micro/go-micro/v3/events"
	"github.com/micro/micro/v3/service/events/client"
)

// DefaultStream is the default events stream implementation
var DefaultStream events.Stream = client.NewStream()

// Publish an event to a topic
func Publish(topic string, msg interface{}, opts ...events.PublishOption) error {
	return DefaultStream.Publish(topic, msg, opts...)
}

// Subscribe to events
func Subscribe(topic string, opts ...events.SubscribeOption) (<-chan events.Event, error) {
	return DefaultStream.Subscribe(topic, opts...)
}
