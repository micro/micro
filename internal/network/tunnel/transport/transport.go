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
// Original source: github.com/micro/go-micro/v3/network/tunnel/mucp/transport/transport.go

// Package transport provides a tunnel transport
package transport

import (
	"context"

	"github.com/micro/micro/v3/internal/network/transport"
	"github.com/micro/micro/v3/internal/network/tunnel"
	"github.com/micro/micro/v3/internal/network/tunnel/mucp"
)

type tunTransport struct {
	options transport.Options

	tunnel tunnel.Tunnel
}

type tunnelKey struct{}

type transportKey struct{}

func (t *tunTransport) Init(opts ...transport.Option) error {
	for _, o := range opts {
		o(&t.options)
	}

	// close the existing tunnel
	if t.tunnel != nil {
		t.tunnel.Close()
	}

	// get the tunnel
	tun, ok := t.options.Context.Value(tunnelKey{}).(tunnel.Tunnel)
	if !ok {
		tun = mucp.NewTunnel()
	}

	// get the transport
	tr, ok := t.options.Context.Value(transportKey{}).(transport.Transport)
	if ok {
		tun.Init(tunnel.Transport(tr))
	}

	// set the tunnel
	t.tunnel = tun

	return nil
}

func (t *tunTransport) Dial(addr string, opts ...transport.DialOption) (transport.Client, error) {
	if err := t.tunnel.Connect(); err != nil {
		return nil, err
	}

	c, err := t.tunnel.Dial(addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (t *tunTransport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	if err := t.tunnel.Connect(); err != nil {
		return nil, err
	}

	l, err := t.tunnel.Listen(addr)
	if err != nil {
		return nil, err
	}

	return &tunListener{l}, nil
}

func (t *tunTransport) Options() transport.Options {
	return t.options
}

func (t *tunTransport) String() string {
	return "tunnel"
}

// NewTransport honours the initialiser used in
func NewTransport(opts ...transport.Option) transport.Transport {
	t := &tunTransport{
		options: transport.Options{},
	}

	// initialise
	t.Init(opts...)

	return t
}

// WithTransport sets the internal tunnel
func WithTunnel(t tunnel.Tunnel) transport.Option {
	return func(o *transport.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tunnelKey{}, t)
	}
}

// WithTransport sets the internal transport
func WithTransport(t transport.Transport) transport.Option {
	return func(o *transport.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, transportKey{}, t)
	}
}
