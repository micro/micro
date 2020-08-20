package util

import (
	"time"

	"github.com/micro/go-micro/v3/events"
	pb "github.com/micro/micro/v3/service/events/proto"
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
