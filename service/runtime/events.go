package runtime

import "github.com/micro/go-micro/v3/runtime"

const (
	// EventServiceCreated is the topic events are published to when a service is created
	EventServiceCreated = "runtime.service.created"
	// EventServiceUpdated is the topic events are published to when a service is updated
	EventServiceUpdated = "runtime.service.updated"
	// EventServiceDeleted is the topic events are published to when a service is deleted
	EventServiceDeleted = "runtime.service.deleted"
)

// EventPayload which is published with runtime events
type EventPayload struct {
	Service   *runtime.Service
	Namespace string
}
