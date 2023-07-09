package handler

import (
	"context"

	pb "micro.dev/v4/proto/events"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/events"
	goevents "micro.dev/v4/service/events"
	"micro.dev/v4/service/events/util"
	"micro.dev/v4/util/auth/namespace"
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
