package handler

import (
	"context"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/broker"
	pb "github.com/micro/go-micro/v2/broker/service/proto"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/internal/namespace"
)

const defaultNamespace = "micro"

type Broker struct {
	Broker broker.Broker
}

func (b *Broker) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.Empty) error {
	ns := namespace.FromContext(ctx)

	// authorize the request
	if err := authorizeNamespace(ctx, ns); err != nil {
		return err
	}

	log.Debugf("Publishing message to %s topic in the %v namespace", req.Topic, ns)
	err := b.Broker.Publish(ns+"."+req.Topic, &broker.Message{
		Header: req.Message.Header,
		Body:   req.Message.Body,
	})
	log.Debugf("Published message to %s topic in the %v namespace", req.Topic, ns)
	if err != nil {
		return errors.InternalServerError("go.micro.broker", err.Error())
	}
	return nil
}

func (b *Broker) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Broker_SubscribeStream) error {
	ns := namespace.FromContext(ctx)

	// authorize the request
	if err := authorizeNamespace(ctx, ns); err != nil {
		return err
	}

	// message handler to stream back messages from broker
	errChan := make(chan error, 1)
	handler := func(p broker.Event) error {
		if err := stream.Send(&pb.Message{
			Header: p.Message().Header,
			Body:   p.Message().Body,
		}); err != nil {
			select {
			case errChan <- err:
				return err
			default:
				return err
			}
		}
		return nil
	}

	log.Debugf("Subscribing to %s topic in namespace %v", req.Topic, ns)
	sub, err := b.Broker.Subscribe(ns+"."+req.Topic, handler, broker.Queue(ns+"."+req.Queue))
	if err != nil {
		return errors.InternalServerError("go.micro.broker", err.Error())
	}
	defer func() {
		log.Debugf("Unsubscribing from topic %s in namespace %v", req.Topic, ns)
		sub.Unsubscribe()
	}()

	select {
	case <-ctx.Done():
		log.Debugf("Context done for subscription to topic %s", req.Topic)
		return nil
	case err := <-errChan:
		log.Debugf("Subscription error for topic %s: %v", req.Topic, err)
		return err
	}
}

// authorizeNamespace returns an error if the context doesn't have access to the given namespace
func authorizeNamespace(ctx context.Context, namespace string) error {
	if namespace == defaultNamespace {
		return nil
	}

	// accounts are always required so we can identify the caller. If auth is not configured, the noop
	// auth implementation will return a blank account with the default domain set, allowing the caller
	// access to all resources
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("go.micro.broker", "An account is required")
	}

	// ensure the account is requesing access to it's own domain
	if acc.Issuer != namespace {
		return errors.Forbidden("go.micro.broker", "An account issued by %v is required", namespace)
	}

	return nil
}
