package client

import (
	pb "github.com/micro/micro/v3/proto/events"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/events/util"
)

// NewStore returns an initialized store handler
func NewStore() events.Store {
	return new(store)
}

type store struct {
	Client pb.StoreService
}

func (s *store) Read(topic string, opts ...events.ReadOption) ([]*events.Event, error) {
	// parse the options
	var options events.ReadOptions
	for _, o := range opts {
		o(&options)
	}

	// execute the RPC
	rsp, err := s.client().Read(context.DefaultContext, &pb.ReadRequest{
		Topic:  topic,
		Limit:  uint64(options.Limit),
		Offset: uint64(options.Offset),
	}, client.WithAuthToken())
	if err != nil {
		return nil, err
	}

	// serialize the response
	result := make([]*events.Event, len(rsp.Events))
	for i, r := range rsp.Events {
		ev := util.DeserializeEvent(r)
		result[i] = &ev
	}

	return result, nil
}

func (s *store) Write(ev *events.Event, opts ...events.WriteOption) error {
	// parse options
	var options events.WriteOptions
	for _, o := range opts {
		o(&options)
	}

	// start the stream
	_, err := s.client().Write(context.DefaultContext, &pb.WriteRequest{
		Event: &pb.Event{
			Id:        ev.ID,
			Topic:     ev.Topic,
			Metadata:  ev.Metadata,
			Payload:   ev.Payload,
			Timestamp: ev.Timestamp.Unix(),
		},
	}, client.WithAuthToken())

	return err
}

// this is a tmp solution since the client isn't initialized when NewStream is called. There is a
// fix in the works in another PR.
func (s *store) client() pb.StoreService {
	if s.Client == nil {
		s.Client = pb.NewStoreService("events", client.DefaultClient)
	}
	return s.Client
}
