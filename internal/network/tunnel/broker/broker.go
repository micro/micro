// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/network/tunnel/broker/broker.go

// Package broker is a tunnel broker
package broker

import (
	"context"

	"github.com/micro/micro/v3/internal/network/transport"
	"github.com/micro/micro/v3/internal/network/tunnel"
	"github.com/micro/micro/v3/internal/network/tunnel/mucp"
	"github.com/micro/micro/v3/service/broker"
)

type tunBroker struct {
	opts   broker.Options
	tunnel tunnel.Tunnel
}

type tunSubscriber struct {
	topic   string
	handler broker.Handler
	opts    broker.SubscribeOptions

	closed   chan bool
	listener tunnel.Listener
}

type tunEvent struct {
	topic   string
	message *broker.Message
}

// used to access tunnel from options context
type tunnelKey struct{}
type tunnelAddr struct{}

func (t *tunBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *tunBroker) Options() broker.Options {
	return t.opts
}

func (t *tunBroker) Address() string {
	return t.tunnel.Address()
}

func (t *tunBroker) Connect() error {
	return t.tunnel.Connect()
}

func (t *tunBroker) Disconnect() error {
	return t.tunnel.Close()
}

func (t *tunBroker) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	// TODO: this is probably inefficient, we might want to just maintain an open connection
	// it may be easier to add broadcast to the tunnel
	c, err := t.tunnel.Dial(topic, tunnel.DialMode(tunnel.Multicast))
	if err != nil {
		return err
	}
	defer c.Close()

	return c.Send(&transport.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (t *tunBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	l, err := t.tunnel.Listen(topic, tunnel.ListenMode(tunnel.Multicast))
	if err != nil {
		return nil, err
	}

	var options broker.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	tunSub := &tunSubscriber{
		topic:    topic,
		handler:  h,
		opts:     options,
		closed:   make(chan bool),
		listener: l,
	}

	// start processing
	go tunSub.run()

	return tunSub, nil
}

func (t *tunBroker) String() string {
	return "tunnel"
}

func (t *tunSubscriber) run() {
	for {
		// accept a new connection
		c, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.closed:
				return
			default:
				continue
			}
		}

		// receive message
		m := new(transport.Message)
		if err := c.Recv(m); err != nil {
			c.Close()
			continue
		}

		// close the connection
		c.Close()

		// handle the message
		go t.handler(&broker.Message{
			Header: m.Header,
			Body:   m.Body,
		})
	}
}

func (t *tunSubscriber) Options() broker.SubscribeOptions {
	return t.opts
}

func (t *tunSubscriber) Topic() string {
	return t.topic
}

func (t *tunSubscriber) Unsubscribe() error {
	select {
	case <-t.closed:
		return nil
	default:
		close(t.closed)
		return t.listener.Close()
	}
}

func (t *tunEvent) Topic() string {
	return t.topic
}

func (t *tunEvent) Message() *broker.Message {
	return t.message
}

func (t *tunEvent) Ack() error {
	return nil
}

func (t *tunEvent) Error() error {
	return nil
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	t, ok := options.Context.Value(tunnelKey{}).(tunnel.Tunnel)
	if !ok {
		t = mucp.NewTunnel()
	}

	a, ok := options.Context.Value(tunnelAddr{}).(string)
	if ok {
		// initialise address
		t.Init(tunnel.Address(a))
	}

	if len(options.Addrs) > 0 {
		// initialise nodes
		t.Init(tunnel.Nodes(options.Addrs...))
	}

	return &tunBroker{
		opts:   options,
		tunnel: t,
	}
}

// WithAddress sets the tunnel address
func WithAddress(a string) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tunnelAddr{}, a)
	}
}

// WithTunnel sets the internal tunnel
func WithTunnel(t tunnel.Tunnel) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tunnelKey{}, t)
	}
}
