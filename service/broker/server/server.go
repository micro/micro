package handler

import (
	"context"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/logger"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v2/internal/namespace"
	"github.com/micro/micro/v2/service"
	mubroker "github.com/micro/micro/v2/service/broker"
	pb "github.com/micro/micro/v2/service/broker/proto"
)

var (
	name    = "go.micro.broker"
	address = ":8001"
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
	mubroker.DefaultBroker.Connect()

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
	if err := namespace.Authorize(ctx, ns); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.broker", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.broker", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.broker", err.Error())
	}

	log.Debugf("Publishing message to %s topic in the %v namespace", req.Topic, ns)
	err := mubroker.DefaultBroker.Publish(ns+"."+req.Topic, &broker.Message{
		Header: req.Message.Header,
		Body:   req.Message.Body,
	})
	log.Debugf("Published message to %s topic in the %v namespace", req.Topic, ns)
	if err != nil {
		return errors.InternalServerError("go.micro.broker", err.Error())
	}
	return nil
}

func (h *handler) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Broker_SubscribeStream) error {
	ns := namespace.FromContext(ctx)
	errChan := make(chan error, 1)

	// authorize the request
	if err := namespace.Authorize(ctx, ns); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.broker", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.broker", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.broker", err.Error())
	}

	// message handler to stream back messages from broker
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
	sub, err := mubroker.DefaultBroker.Subscribe(ns+"."+req.Topic, handler, broker.Queue(ns+"."+req.Queue))
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
