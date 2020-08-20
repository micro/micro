package events

import (
	"github.com/micro/go-micro/v3/events"
	"github.com/micro/micro/v3/service/events/client"
)

var (
	// DefaultStream is the default events stream implementation
	DefaultStream events.Stream = client.NewStream()
	// DefaultStore is the default events store implementation
	DefaultStore events.Store = client.NewStore()
)

// Publish an event to a topic
func Publish(topic string, msg interface{}, opts ...events.PublishOption) error {
	return DefaultStream.Publish(topic, msg, opts...)
}

// Subscribe to events
func Subscribe(topic string, opts ...events.SubscribeOption) (<-chan events.Event, error) {
	return DefaultStream.Subscribe(topic, opts...)
}

// Read events for a topic
func Read(topic string, opts ...events.ReadOption) ([]*events.Event, error) {
	return DefaultStore.Read(topic, opts...)
}
