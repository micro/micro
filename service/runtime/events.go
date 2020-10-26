package runtime

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
	EventResourceQuotaCreated = "resourcequota.created"
	EventResourceQuotaUpdated = "resourcequota.updated"
	EventResourceQuotaDeleted = "resourcequota.deleted"
)

// EventPayload which is published with runtime events
type EventPayload struct {
	Type      string
	Service   *Service
	Namespace string
}

// EventResourcePayload which is published with runtime resource events
type EventResourcePayload struct {
	Type          string
	Name          string
	Namespace     string
	NetworkPolicy *NetworkPolicy
	ResourceQuota *ResourceQuota
	Service       *Service
}
