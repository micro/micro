package runtime

import "github.com/micro/go-micro/v3/runtime"

const (
	// EventTopic the events are published to
	EventTopic = "runtime"

	// EventServiceCreated is the topic events are published to when a service is created
	EventServiceCreated = "service.created"
	// EventServiceUpdated is the topic events are published to when a service is updated
	EventServiceUpdated = "service.updated"
	// EventServiceDeleted is the topic events are published to when a service is deleted
	EventServiceDeleted = "service.deleted"
)

// EventPayload which is published with runtime events
type EventPayload struct {
	Type      string
	Service   *runtime.Service
	Namespace string
}
