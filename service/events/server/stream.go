package server

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/events"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/events/util"
	"github.com/micro/micro/v3/service/logger"
)

type Stream struct{}

func (s *Stream) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
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
		return errors.BadRequest("events.Stream.Publish", events.ErrMissingTopic.Error())
	}

	// parse options
	var opts []events.PublishOption
	if req.Timestamp > 0 {
		opts = append(opts, events.WithTimestamp(time.Unix(req.Timestamp, 0)))
	}
	if req.Metadata != nil {
		opts = append(opts, events.WithMetadata(req.Metadata))
	}

	// publish the event
	if err := events.Publish(req.Topic, req.Payload, opts...); err != nil {
		return errors.InternalServerError("events.Stream.Publish", err.Error())
	}

	// write the event to the store
	event := events.Event{
		ID:        uuid.New().String(),
		Metadata:  req.Metadata,
		Payload:   req.Payload,
		Topic:     req.Topic,
		Timestamp: time.Unix(req.Timestamp, 0),
	}
	if err := events.DefaultStore.Write(&event, events.WithTTL(time.Hour*24)); err != nil {
		logger.Errorf("Error writing event %v to store: %v", event.ID, err)
	}

	return nil
}

func (s *Stream) Consume(ctx context.Context, req *pb.ConsumeRequest, rsp pb.Stream_ConsumeStream) error {
	// authorize the request
	if err := namespace.Authorize(ctx, namespace.DefaultNamespace); err == namespace.ErrForbidden {
		return errors.Forbidden("events.Stream.Publish", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("events.Stream.Publish", err.Error())
	} else if err != nil {
		return errors.InternalServerError("events.Stream.Publish", err.Error())
	}

	// parse options
	opts := []events.ConsumeOption{}
	if req.Offset > 0 {
		opts = append(opts, events.WithOffset(time.Unix(req.Offset, 0)))
	}
	if len(req.Group) > 0 {
		opts = append(opts, events.WithGroup(req.Group))
	}
	if !req.AutoAck {
		opts = append(opts, events.WithAutoAck(req.AutoAck, time.Duration(req.AckWait)/time.Nanosecond))
	}
	if req.RetryLimit > -1 {
		opts = append(opts, events.WithRetryLimit(int(req.RetryLimit)))
	}

	// create the subscriber
	evChan, err := events.Consume(req.Topic, opts...)
	if err != nil {
		return errors.InternalServerError("events.Stream.Consume", err.Error())
	}

	type eventSent struct {
		sent  time.Time
		event events.Event
	}
	ackMap := map[string]eventSent{}
	mutex := sync.RWMutex{}
	go func() {
		for {
			req := pb.AckRequest{}
			if err := rsp.RecvMsg(&req); err != nil {
				return
			}
			mutex.RLock()
			ev, ok := ackMap[req.Id]
			mutex.RUnlock()
			if !ok {
				// not found, probably timed out after ackWait
				continue
			}
			if req.Success {
				ev.event.Ack()
			} else {
				ev.event.Nack()
			}
			mutex.Lock()
			delete(ackMap, req.Id)
			mutex.Unlock()
		}
	}()
	for {
		// Do any clean up of ackMap where we haven't got a response
		now := time.Now()
		for k, v := range ackMap {
			if v.sent.Add(2 * time.Duration(req.AckWait)).Before(now) {
				mutex.Lock()
				delete(ackMap, k)
				mutex.Unlock()
			}
		}
		ev, ok := <-evChan
		if !ok {
			rsp.Close()
			return nil
		}
		if !req.AutoAck {
			// track the acks
			mutex.Lock()
			ackMap[ev.ID] = eventSent{event: ev, sent: time.Now()}
			mutex.Unlock()
		}

		if err := rsp.Send(util.SerializeEvent(&ev)); err != nil {
			return err
		}
	}
}
