package broker

import (
	"context"
	"time"

	authns "github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/broker"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
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
	// authorize the request
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("broker.Broker.Publish", authns.ErrForbidden.Error())
	}

	// validate the request
	if req.Message == nil {
		return errors.BadRequest("broker.Broker.Publish", "Missing message")
	}

	// ensure the header is not nil
	if req.Message.Header == nil {
		req.Message.Header = map[string]string{}
	}

	// set any headers which aren't already set
	if md, ok := metadata.FromContext(ctx); ok {
		for k, v := range md {
			if _, ok := req.Message.Header[k]; !ok {
				req.Message.Header[k] = v
			}
		}
	}

	logger.Debugf("Publishing message to %s topic in the %v namespace", req.Topic, acc.Issuer)
	err := broker.DefaultBroker.Publish(acc.Issuer+"."+req.Topic, &broker.Message{
		Header: req.Message.Header,
		Body:   req.Message.Body,
	})
	logger.Debugf("Published message to %s topic in the %v namespace", req.Topic, acc.Issuer)
	if err != nil {
		return errors.InternalServerError("broker.Broker.Publish", err.Error())
	}
	return nil
}

func (h *handler) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Broker_SubscribeStream) error {
	// authorize the request
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("broker.Broker.Subscribe", authns.ErrForbidden.Error())
	}
	ns := acc.Issuer

	errChan := make(chan error, 1)

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

	logger.Debugf("Subscribing to %s topic in namespace %v", req.Topic, ns)
	sub, err := broker.DefaultBroker.Subscribe(ns+"."+req.Topic, handler, broker.Queue(ns+"."+req.Queue))
	if err != nil {
		return errors.InternalServerError("broker.Broker.Subscribe", err.Error())
	}
	defer func() {
		logger.Debugf("Unsubscribing from topic %s in namespace %v", req.Topic, ns)
		sub.Unsubscribe()
	}()

	select {
	case <-ctx.Done():
		logger.Debugf("Context done for subscription to topic %s", req.Topic)
		return nil
	case err := <-errChan:
		logger.Debugf("Subscription error for topic %s: %v", req.Topic, err)
		return err
	}
}
