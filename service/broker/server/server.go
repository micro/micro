package handler

import (
	"context"
	"time"

	authns "github.com/micro/micro/v3/internal/auth/namespace"
	"github.com/micro/micro/v3/internal/namespace"
	pb "github.com/micro/micro/v3/proto/broker"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"
)

var (
	name    = "broker"
	address = ":8003"
)

// Run the micro broker
func Run(ctx *cli.Context) error {
	srvOpts := []service.Option{
		service.Name(name),
		service.Address(address),
	}

	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, service.RegisterInterval(i*time.Second))
	}

	// new service
	srv := service.New(srvOpts...)

	// connect to the broker
	broker.DefaultBroker.Connect()

	// register the broker handler
	pb.RegisterBrokerHandler(srv.Server(), new(handler))

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
	return nil
}

type handler struct{}

func (h *handler) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.Empty) error {
	ns := namespace.FromContext(ctx)

	// authorize the request
	if err := authns.Authorize(ctx, ns); err == authns.ErrForbidden {
		return errors.Forbidden("broker.Broker.Publish", err.Error())
	} else if err == authns.ErrUnauthorized {
		return errors.Unauthorized("broker.Broker.Publish", err.Error())
	} else if err != nil {
		return errors.InternalServerError("broker.Broker.Publish", err.Error())
	}

	log.Debugf("Publishing message to %s topic in the %v namespace", req.Topic, ns)
	err := broker.DefaultBroker.Publish(ns+"."+req.Topic, &broker.Message{
		Header: req.Message.Header,
		Body:   req.Message.Body,
	})
	log.Debugf("Published message to %s topic in the %v namespace", req.Topic, ns)
	if err != nil {
		return errors.InternalServerError("broker.Broker.Publish", err.Error())
	}
	return nil
}

func (h *handler) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Broker_SubscribeStream) error {
	ns := namespace.FromContext(ctx)
	errChan := make(chan error, 1)

	// authorize the request
	if err := authns.Authorize(ctx, ns); err == authns.ErrForbidden {
		return errors.Forbidden("broker.Broker.Subscribe", err.Error())
	} else if err == authns.ErrUnauthorized {
		return errors.Unauthorized("broker.Broker.Subscribe", err.Error())
	} else if err != nil {
		return errors.InternalServerError("broker.Broker.Subscribe", err.Error())
	}

	// message handler to stream back messages from broker
	handler := func(m *broker.Message) error {
		if err := stream.Send(&pb.Message{
			Header: m.Header,
			Body:   m.Body,
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
	sub, err := broker.DefaultBroker.Subscribe(ns+"."+req.Topic, handler, broker.Queue(ns+"."+req.Queue))
	if err != nil {
		return errors.InternalServerError("broker.Broker.Subscribe", err.Error())
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
