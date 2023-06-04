package transport

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/gob"
	"errors"
	"net"
	"time"

	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/network/transport"
	maddr "github.com/micro/micro/v3/util/addr"
	mnet "github.com/micro/micro/v3/util/net"
	mls "github.com/micro/micro/v3/util/tls"
	"tailscale.com/tsnet"
)

type tsTransport struct {
	opts transport.Options
	srv  *tsnet.Server
}

type tsTransportClient struct {
	dialOpts transport.DialOptions
	conn     net.Conn
	enc      *gob.Encoder
	dec      *gob.Decoder
	encBuf   *bufio.Writer
	timeout  time.Duration
}

type tsTransportSocket struct {
	conn    net.Conn
	enc     *gob.Encoder
	dec     *gob.Decoder
	encBuf  *bufio.Writer
	timeout time.Duration
}

type tsTransportListener struct {
	listener net.Listener
	timeout  time.Duration
}

func (t *tsTransportClient) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tsTransportClient) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tsTransportClient) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	if err := t.enc.Encode(m); err != nil {
		return err
	}
	return t.encBuf.Flush()
}

func (t *tsTransportClient) Recv(m *transport.Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	return t.dec.Decode(&m)
}

func (t *tsTransportClient) Close() error {
	return t.conn.Close()
}

func (t *tsTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tsTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tsTransportSocket) Recv(m *transport.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}

	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	return t.dec.Decode(&m)
}

func (t *tsTransportSocket) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	if err := t.enc.Encode(m); err != nil {
		return err
	}
	return t.encBuf.Flush()
}

func (t *tsTransportSocket) Close() error {
	return t.conn.Close()
}

func (t *tsTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *tsTransportListener) Close() error {
	return t.listener.Close()
}

func (t *tsTransportListener) Accept(fn func(transport.Socket)) error {
	var tempDelay time.Duration

	for {
		c, err := t.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Errorf("http: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		encBuf := bufio.NewWriter(c)
		sock := &tsTransportSocket{
			timeout: t.timeout,
			conn:    c,
			encBuf:  encBuf,
			enc:     gob.NewEncoder(encBuf),
			dec:     gob.NewDecoder(c),
		}

		go func() {
			// TODO: think of a better error response strategy
			defer func() {
				if r := recover(); r != nil {
					sock.Close()
				}
			}()

			fn(sock)
		}()
	}
}

func (t *tsTransport) Dial(addr string, opts ...transport.DialOption) (transport.Client, error) {
	dopts := transport.DialOptions{
		Timeout: transport.DefaultDialTimeout,
	}

	for _, opt := range opts {
		opt(&dopts)
	}

	var conn net.Conn
	var err error

	// TODO: support dial option here rather than using internal config
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		// TODO: dopts.Timeout
		conn, err = t.srv.Dial(context.TODO(), "tcp", addr)
		if err != nil {
			return nil, err
		}
		conn = tls.Client(conn, config)
	} else {
		// TODO: dopts.Timeout
		conn, err = t.srv.Dial(context.TODO(), "tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	encBuf := bufio.NewWriter(conn)

	return &tsTransportClient{
		dialOpts: dopts,
		conn:     conn,
		encBuf:   encBuf,
		enc:      gob.NewEncoder(encBuf),
		dec:      gob.NewDecoder(conn),
		timeout:  t.opts.Timeout,
	}, nil
}

func (t *tsTransport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	var options transport.ListenOptions
	for _, o := range opts {
		o(&options)
	}

	var l net.Listener
	var err error

	// TODO: support use of listen options
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig

		fn := func(addr string) (net.Listener, error) {
			if config == nil {
				hosts := []string{addr}

				// check if its a valid host:port
				if host, _, err := net.SplitHostPort(addr); err == nil {
					if len(host) == 0 {
						hosts = maddr.IPs()
					} else {
						hosts = []string{host}
					}
				}

				// generate a certificate
				cert, err := mls.Certificate(hosts...)
				if err != nil {
					return nil, err
				}
				config = &tls.Config{Certificates: []tls.Certificate{cert}}
			}
			ln, err := t.srv.Listen("tcp", addr)
			if err != nil {
				return nil, err
			}
			// TODO: tailscale.GetCertificate
			return tls.NewListener(ln, config), err
		}

		l, err = mnet.Listen(addr, fn)
	} else {
		fn := func(addr string) (net.Listener, error) {
			return t.srv.Listen("tcp", addr)
		}

		l, err = mnet.Listen(addr, fn)
	}

	if err != nil {
		return nil, err
	}

	return &tsTransportListener{
		timeout:  t.opts.Timeout,
		listener: l,
	}, nil
}

func (t *tsTransport) Init(opts ...transport.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *tsTransport) Options() transport.Options {
	return t.opts
}

func (t *tsTransport) String() string {
	return "tailscale"
}

func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &tsTransport{
		opts: options,
		srv:  new(tsnet.Server),
	}
}
