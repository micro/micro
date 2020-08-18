package server

import (
	"context"
	"time"

	"github.com/micro/cli/v2"
	goevents "github.com/micro/go-micro/v3/events"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	pb "github.com/micro/micro/v3/service/events/proto"
	"github.com/micro/micro/v3/service/logger"
)

// Run the micro broker
func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("events"),
	)

	// register the broker handler
	pb.RegisterStreamHandler(srv.Server(), new(handler))

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	return nil
}

type handler struct{}

func (h *handler) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
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
	if req.Payload != nil {
		opts = append(opts, goevents.WithPayload(req.Payload))
	}
	if req.Metadata != nil {
		opts = append(opts, goevents.WithMetadata(req.Metadata))
	}

	// publish the event
	if err := events.Publish(req.Topic, opts...); err != nil {
		return errors.InternalServerError("events.Stream.Publish", err.Error())
	}

	return nil
}

func (h *handler) Subscribe(ctx context.Context, req *pb.SubscribeRequest, rsp pb.Stream_SubscribeStream) error {
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
	if len(req.Topic) > 0 {
		opts = append(opts, goevents.WithTopic(req.Topic))
	}

	// create the subscriber
	evChan, err := events.Subscribe(opts...)
	if err != nil {
		return errors.InternalServerError("events.Stream.Subscribe", err.Error())
	}

	go func() {
		for {
			ev, ok := <-evChan
			if !ok {
				rsp.Close()
				return
			}

			// todo: handle error, don't ack the event
			rsp.Send(&pb.Event{
				Id:        ev.ID,
				Topic:     ev.Topic,
				Metadata:  ev.Metadata,
				Payload:   ev.Payload,
				Timestamp: ev.Timestamp.Unix(),
			})
		}
	}()

	return nil
}
