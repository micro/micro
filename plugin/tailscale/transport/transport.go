package transport

import (
	"context"
	"net"

	"github.com/micro/micro/v3/service/network/transport"
	"github.com/micro/micro/v3/service/network/transport/http"
	"tailscale.com/tsnet"
)

func NewTransport(opts ...transport.Option) transport.Transport {
	tr := http.NewTransport(opts...)
	srv := new(tsnet.Server)

	return &tailscale{
		tr:  tr,
		srv: srv,
	}
}

type tailscale struct {
	tr transport.Transport

	srv *tsnet.Server
}

func (t *tailscale) Init(opts ...transport.Option) error {
	return t.tr.Init(opts...)
}

func (t *tailscale) Options() transport.Options {
	return t.tr.Options()
}

func (t *tailscale) Dial(addr string, opts ...transport.DialOption) (transport.Client, error) {
	tsDialer := func(addr string) (net.Conn, error) {
		return t.srv.Dial(context.TODO(), "tcp", addr)
	}
	opts = append(opts, transport.WithDialFunc(tsDialer))
	return t.tr.Dial(addr, opts...)
}

func (t *tailscale) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	l, err := t.srv.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	opts = append(opts, http.NetListener(l))
	return t.tr.Listen(addr, opts...)
}

func (t *tailscale) String() string {
	return "tailscale"
}
