package util

import (
	"time"

	pb "micro.dev/v4/proto/events"
	"micro.dev/v4/service/events"
)

func SerializeEvent(ev *events.Event) *pb.Event {
	return &pb.Event{
		Id:        ev.ID,
		Topic:     ev.Topic,
		Metadata:  ev.Metadata,
		Payload:   ev.Payload,
		Timestamp: ev.Timestamp.Unix(),
	}
}

func DeserializeEvent(ev *pb.Event) events.Event {
	return events.Event{
		ID:        ev.Id,
		Topic:     ev.Topic,
		Metadata:  ev.Metadata,
		Payload:   ev.Payload,
		Timestamp: time.Unix(ev.Timestamp, 0),
	}
}
