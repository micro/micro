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
	EventServiceDeleted       = "service.deleted"
	EventNamespaceCreated     = "namespace.created"
	EventNamespaceDeleted     = "namespace.deleted"
	EventNetworkPolicyCreated = "networkpolicy.created"
	EventNetworkPolicyUpdated = "networkpolicy.updated"
	EventNetworkPolicyDeleted = "networkpolicy.deleted"
)

// EventPayload which is published with runtime events
type EventPayload struct {
	Type      string
	Service   *runtime.Service
	Namespace string
}

// EventNamespacePayload which is published with runtime namespace events
type EventNamespacePayload struct {
	Type      string
	Namespace string
}

// EventNetworkPolicyPayload which is published with runtime networkpolicy events
type EventNetworkPolicyPayload struct {
	Type      string
	Name      string
	Namespace string
}
