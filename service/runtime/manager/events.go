package manager

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/store"
)

type eventType string

const (
	eventTypeCreated eventType = "created"
	eventTypeUpdated eventType = "updated"
	eventTypeDeleted eventType = "deleted"
)

var (
	// eventTTL is the duration events will perist in the store before expiring
	eventTTL = time.Minute * 10
	// eventPollFrequency is the max frequency the manager will check for new events in the store
	eventPollFrequency = time.Minute
)

// eventPrefix is prefixed to the key for event records
const eventPrefix = "event/"

// event is the object written to the store
type event struct {
	ID      string                 `json:"id"`
	Type    eventType              `json:"type"`
	Service *runtime.Service       `json:"service"`
	Options *runtime.CreateOptions `json:"options"`
}

// key the event should be written to in the store
func (e *event) Key() string {
	return eventPrefix + e.ID
}

// publishEvent will write the event to the global store and immediately process the event
func (m *manager) publishEvent(eType eventType, srv *runtime.Service, opts *runtime.CreateOptions) error {
	e := &event{
		ID:      uuid.New().String(),
		Type:    eType,
		Service: srv,
		Options: opts,
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	record := &store.Record{
		Key:    e.Key(),
		Value:  bytes,
		Expiry: eventTTL,
	}

	if err := m.options.Store.Write(record); err != nil {
		return err
	}

	go m.processEvent(e.ID)
	return nil
}

// watchEvents polls the store for events periodically and processes
// them if they have not already done so
func (m *manager) watchEvents() {
	ticker := time.NewTicker(eventPollFrequency)

	for {
		<-ticker.C

		// get the keys of the events
		eventKeys, err := m.options.Store.List(store.ListPrefix(eventPrefix))
		if err != nil {
			logger.Warn("Error listing events: %v", err)
			continue
		}

		// loop through every event
		for _, key := range eventKeys {
			m.processEvent(strings.TrimPrefix(key, eventPrefix))
		}
	}
}

// processEvent will take an event key, verify it hasn't been consumed
// and then execute it
func (m *manager) processEvent(id string) {
	// check to see if the event has been processed before
	if _, err := m.eventsConsumed.Read(id); err != store.ErrNotFound {
		return
	}

	// TODO: process the event

	// write to the store indicating the event has been consumed. We double
	// the ttl to safely know the event will expire before this record
	m.eventsConsumed.Write(&store.Record{Key: id, Expiry: eventTTL * 2})
}
