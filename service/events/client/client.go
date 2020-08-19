package client

import (
	"encoding/json"
	"time"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/events"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	pb "github.com/micro/micro/v3/service/events/proto"
)

// NewStream returns an initialized stream service
func NewStream() events.Stream {
	return new(stream)
}

type stream struct {
	Client pb.StreamService
}

func (s *stream) Publish(topic string, opts ...events.PublishOption) error {
	// parse the options
	options := events.PublishOptions{
		Timestamp: time.Now(),
	}
	for _, o := range opts {
		o(&options)
	}

	// encode the message if it's not already encoded
	var payload []byte
	if p, ok := options.Payload.([]byte); ok {
		payload = p
	} else {
		p, err := json.Marshal(options.Payload)
		if err != nil {
			return events.ErrEncodingMessage
		}
		payload = p
	}

	// execute the RPC
	_, err := s.client().Publish(context.DefaultContext, &pb.PublishRequest{
		Topic:     topic,
		Payload:   payload,
		Metadata:  options.Metadata,
		Timestamp: options.Timestamp.Unix(),
	}, goclient.WithAuthToken())

	return err
}

func (s *stream) Subscribe(opts ...events.SubscribeOption) (<-chan events.Event, error) {
	// parse options
	var options events.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	// start the stream
	stream, err := s.client().Subscribe(context.DefaultContext, &pb.SubscribeRequest{
		Queue:       options.Queue,
		Topic:       options.Topic,
		StartAtTime: options.StartAtTime.Unix(),
	}, goclient.WithAuthToken())
	if err != nil {
		return nil, err
	}

	evChan := make(chan events.Event)
	go func() {
		for {
			ev, err := stream.Recv()
			if err != nil {
				close(evChan)
				return
			}

			evChan <- events.Event{
				ID:        ev.Id,
				Topic:     ev.Topic,
				Metadata:  ev.Metadata,
				Payload:   ev.Payload,
				Timestamp: time.Unix(ev.Timestamp, 0),
			}
		}
	}()

	return evChan, nil
}

// this is a tmp solution since the client isn't initialized when NewStream is called. There is a
// fix in the works in another PR.
func (s *stream) client() pb.StreamService {
	if s.Client == nil {
		s.Client = pb.NewStreamService("events", client.DefaultClient)
	}
	return s.Client
}
