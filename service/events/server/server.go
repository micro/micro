package server

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/micro/cli"
	events "github.com/micro/go-micro/v3/events"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/micro/v3/service/events/proto"
)

// Events handles RPC requests
type Events struct {
	subs []subscriber
	sync.RWMutex
}

// subscriber is an entity which has subscibed to a topic
type subscriber struct {
	Topic  string
	Queue  string
	Prefix bool
	Stream pb.Broker_SubscribeStream
}

// ShouldSend returns a boolean indicating if an event should be sent to the subscriber
func (s *subscriber) ShouldSend(e *events.Event) bool {
	if s.Prefix {
		return strings.HasPrefix(s.Topic, e.Topic)
	}
	return s.Topic == e.Topic
}

// Publish an event
func (e *Events) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
	// validate the request
	if len(req.Topic) == 0 {
		return errors.BadRequest("events.Broker.Publish", "Missing topic")
	}

	// construct the event
	ev := &events.Event{
		ID:       uuid.New().String(),
		Topic:    req.Topic,
		Metadata: req.GetOptions().GetMetadata(),
		Payload:  req.GetOptions().GetPayload(),
	}

	if unix := req.GetOptions().GetTimestamp(); unix > 0 {
		ev.Timestamp = time.Unix(unix, 0)
	} else {
		ev.Timestamp = time.Now()
	}

	// todo: write the event to the store

	// proto encode the message so we can send it to the subscribers
	protoEv := &pb.Event{
		Id:       ev.ID,
		Topic:    ev.Topic,
		Metadata: ev.Metadata,
		Payload:  ev.Payload,
	}

	// lock before we iterate over the subscribers
	e.RLock()
	defer e.Unlock()

	// send the message to any subscribers who want the message. Record the name of the queue we send
	// the message to since queues with the same name should only recieve a message once
	var queues []string
	for _, s := range e.subs {
		if contains(queues, s.Queue) {
			continue
		}

		if s.ShouldSend(ev) {
			// todo: handle the error returned when publishing to a stream
			s.Stream.Send(protoEv)
		}

		queues = append(queues, s.Queue)
	}

	return nil
}

// Subscribe to a topic
func (e *Events) Subscribe(ctx context.Context, req *pb.SubscribeRequest, rsp pb.Broker_SubscribeStream) error {
	// validate the request
	if len(req.Topic) == 0 && !req.GetOptions().Prefix {
		return errors.BadRequest("events.Broker.Subscribe", "Topic is required unless the prefix option is specified")
	}

	// construct the subscriber
	sub := subscriber{
		Topic:  req.Topic,
		Prefix: req.GetOptions().Prefix,
		Queue:  req.GetOptions().GetQueue(),
		Stream: rsp,
	}

	// default the queue if none was provided
	if len(sub.Queue) == 0 {
		sub.Queue = uuid.New().String()
	}

	// register the subscriber
	e.Lock()
	defer e.Unlock()
	e.subs = append(e.subs, sub)

	return nil
}

// contains is a helper function which returns a boolean indicating if the value exists in the slice
func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("events"),
		service.Version("latest"),
	)

	// register the handler
	pb.RegisterBrokerHandler(srv.Server(), &Events{
		subs: make([]subscriber, 0),
	})

	// run the service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
