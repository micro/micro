package events

import (
	"context"

	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/events"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	goevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/events/util"
)

type Store struct{}

func (s *Store) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, namespace.DefaultNamespace, "events.Store.Read"); err != nil {
		return err
	}

	// validate the request
	if len(req.Topic) == 0 {
		return errors.BadRequest("events.Store.Read", goevents.ErrMissingTopic.Error())
	}

	// parse options
	var opts []goevents.ReadOption
	if req.Limit > 0 {
		opts = append(opts, goevents.ReadLimit(uint(req.Limit)))
	}
	if req.Offset > 0 {
		opts = append(opts, goevents.ReadOffset(uint(req.Offset)))
	}

	// read from the store
	result, err := events.DefaultStore.Read(req.Topic, opts...)
	if err != nil {
		return errors.InternalServerError("events.Store.Read", err.Error())
	}

	// serialize the result
	rsp.Events = make([]*pb.Event, len(result))
	for i, r := range result {
		rsp.Events[i] = util.SerializeEvent(r)
	}

	return nil
}

func (s *Store) Write(ctx context.Context, req *pb.WriteRequest, rsp *pb.WriteResponse) error {
	return errors.NotImplemented("events.Store.Write", "Writing to the store directly is not supported")
}
