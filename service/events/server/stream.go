package server

import (
	"context"
	"time"

	goevents "github.com/micro/go-micro/v3/events"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	pb "github.com/micro/micro/v3/service/events/proto"
	"github.com/micro/micro/v3/service/events/util"
)

type evStream struct{}

func (s *evStream) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
	// authorize the request
	if err := namespace.Authorize(ctx, namespace.DefaultNamespace); err == namespace.ErrForbidden {
		return errors.Forbidden("events.Stream.Publish", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("events.Stream.Publish", err.Error())
	} else if err != nil {
		return errors.InternalServerError("events.Stream.Publish", err.Error())
	}

	// validate the request
	if len(req.Topic) == 0 {
		return errors.BadRequest("events.Stream.Publish", goevents.ErrMissingTopic.Error())
	}

	// parse options
	var opts []goevents.PublishOption
	if req.Timestamp > 0 {
		opts = append(opts, goevents.WithTimestamp(time.Unix(req.Timestamp, 0)))
	}
	if req.Metadata != nil {
		opts = append(opts, goevents.WithMetadata(req.Metadata))
	}

	// publish the event
	if err := events.Publish(req.Topic, req.Payload, opts...); err != nil {
		return errors.InternalServerError("events.Stream.Publish", err.Error())
	}

	return nil
}

func (s *evStream) Subscribe(ctx context.Context, req *pb.SubscribeRequest, rsp pb.Stream_SubscribeStream) error {
	// authorize the request
	if err := namespace.Authorize(ctx, namespace.DefaultNamespace); err == namespace.ErrForbidden {
		return errors.Forbidden("events.Stream.Publish", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("events.Stream.Publish", err.Error())
	} else if err != nil {
		return errors.InternalServerError("events.Stream.Publish", err.Error())
	}

	// parse options
	var opts []goevents.SubscribeOption
	if req.StartAtTime > 0 {
		opts = append(opts, goevents.WithStartAtTime(time.Unix(req.StartAtTime, 0)))
	}
	if len(req.Queue) > 0 {
		opts = append(opts, goevents.WithQueue(req.Queue))
	}

	// create the subscriber
	evChan, err := events.Subscribe(req.Topic, opts...)
	if err != nil {
		return errors.InternalServerError("events.Stream.Subscribe", err.Error())
	}

	for {
		ev, ok := <-evChan
		if !ok {
			rsp.Close()
			return nil
		}

		if err := rsp.Send(util.SerializeEvent(&ev)); err != nil {
			return err
		}
	}
}
