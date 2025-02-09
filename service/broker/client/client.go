package client

import (
	"time"

	pb "github.com/micro/micro/v5/proto/broker"
	"github.com/micro/micro/v5/service/broker"
	"github.com/micro/micro/v5/service/client"
	"github.com/micro/micro/v5/service/context"
	"github.com/micro/micro/v5/service/logger"
)

var (
	name    = "broker"
	address = ":8003"
)

type serviceBroker struct {
	Addrs   []string
	Client  pb.BrokerService
	options broker.Options
}

func (b *serviceBroker) Address() string {
	return b.Addrs[0]
}

func (b *serviceBroker) Connect() error {
	return nil
}

func (b *serviceBroker) Disconnect() error {
	return nil
}

func (b *serviceBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&b.options)
	}
	b.Client = pb.NewBrokerService(name, client.DefaultClient)
	return nil
}

func (b *serviceBroker) Options() broker.Options {
	return b.options
}

func (b *serviceBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	if logger.V(logger.DebugLevel, logger.DefaultLogger) {
		logger.Debugf("Publishing to topic %s broker %v", topic, b.Addrs)
	}
	_, err := b.Client.Publish(context.DefaultContext, &pb.PublishRequest{
		Topic: topic,
		Message: &pb.Message{
			Header: msg.Header,
			Body:   msg.Body,
		},
	}, client.WithAuthToken(), client.WithAddress(b.Addrs...))
	return err
}

func (b *serviceBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	var options broker.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}
	if logger.V(logger.DebugLevel, logger.DefaultLogger) {
		logger.Debugf("Subscribing to topic %s broker %v", topic, b.Addrs)
	}
	stream, err := b.Client.Subscribe(context.DefaultContext, &pb.SubscribeRequest{
		Topic: topic,
	}, client.WithAuthToken(), client.WithAddress(b.Addrs...), client.WithRequestTimeout(time.Hour))
	if err != nil {
		return nil, err
	}

	sub := &serviceSub{
		topic:   topic,
		handler: handler,
		stream:  stream,
		closed:  make(chan bool),
		options: options,
	}

	go func() {
		for {
			select {
			case <-sub.closed:
				if logger.V(logger.DebugLevel, logger.DefaultLogger) {
					logger.Debugf("Unsubscribed from topic %s", topic)
				}
				return
			default:
				if logger.V(logger.DebugLevel, logger.DefaultLogger) {
					// run the subscriber
					logger.Debugf("Streaming from broker %v to topic [%s]", b.Addrs, topic)
				}
				if err := sub.run(); err != nil {
					if logger.V(logger.DebugLevel, logger.DefaultLogger) {
						logger.Debugf("Resubscribing to topic %s broker %v", topic, b.Addrs)
					}
					stream, err := b.Client.Subscribe(context.DefaultContext, &pb.SubscribeRequest{
						Topic: topic,
					}, client.WithAuthToken(), client.WithAddress(b.Addrs...), client.WithRequestTimeout(time.Hour))
					if err != nil {
						if logger.V(logger.DebugLevel, logger.DefaultLogger) {
							logger.Debugf("Failed to resubscribe to topic %s: %v", topic, err)
						}
						time.Sleep(time.Second)
						continue
					}
					// new stream
					sub.stream = stream
				}
			}
		}
	}()

	return sub, nil
}

func (b *serviceBroker) String() string {
	return "service"
}

func NewBroker(opts ...broker.Option) broker.Broker {
	var options broker.Options
	for _, o := range opts {
		o(&options)
	}

	addrs := options.Addrs
	if len(addrs) == 0 {
		addrs = []string{address}
	}

	return &serviceBroker{
		Addrs:   addrs,
		Client:  pb.NewBrokerService(name, client.DefaultClient),
		options: options,
	}
}
