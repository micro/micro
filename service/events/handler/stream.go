package handler

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	pb "github.com/micro/micro/v3/proto/events"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/events/util"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/util/auth/namespace"
)

type Stream struct{}

func (s *Stream) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, namespace.DefaultNamespace, "events.Stream.Publish"); err != nil {
		return err
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
	if err := namespace.AuthorizeAdmin(ctx, namespace.DefaultNamespace, "events.Stream.Consume"); err != nil {
		return err
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

	// append the context
	opts = append(opts, events.WithContext(ctx))

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
	recvErrChan := make(chan error)
	sendErrChan := make(chan error)

	go func() {
		// process messages from the consumer (probably just ACK messages
		defer close(recvErrChan)
		for {
			select {
			case <-ctx.Done():
				return
			case <-sendErrChan:
				return
			default:
			}

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

	go func() {
		// process messages coming from the stream
		defer close(sendErrChan)
		for {
			// Do any clean up of ackMap where we haven't got a response
			now := time.Now()

			mutex.Lock()
			for k, v := range ackMap {
				if v.sent.Add(2 * time.Duration(req.AckWait)).Before(now) {
					delete(ackMap, k)
				}
			}
			mutex.Unlock()
			var ev events.Event
			var ok bool
			select {
			case <-recvErrChan:
			case <-rsp.Context().Done():
			case <-ctx.Done():
			case ev, ok = <-evChan:
			}
			if !ok {
				return
			}
			if len(ev.ID) == 0 {
				// ignore
				continue
			}
			if !req.AutoAck {
				// track the acks
				mutex.Lock()
				ackMap[ev.ID] = eventSent{event: ev, sent: time.Now()}
				mutex.Unlock()
			}
			if err := rsp.Send(util.SerializeEvent(&ev)); err != nil {
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
	case <-recvErrChan:
	case <-sendErrChan:
	case <-rsp.Context().Done():
	}
	return nil

}
