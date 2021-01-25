// Copyright 2020 Asim Aslam
//
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
// Original source: github.com/micro/go-micro/v3/server/grpc/options.go

package grpc

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/service/broker/memory"
	"github.com/micro/micro/v3/service/registry/mdns"
	"github.com/micro/micro/v3/service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

type codecsKey struct{}
type grpcOptions struct{}
type netListener struct{}
type maxMsgSizeKey struct{}
type maxRecvMsgSizeKey struct{}
type maxSendMsgSizeKey struct{}
type maxConnKey struct{}
type tlsAuth struct{}
type grpcWebOptions struct{}
type grpcWebPort struct{}

// gRPC Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c encoding.Codec) server.Option {
	return func(o *server.Options) {
		codecs := make(map[string]encoding.Codec)
		if o.Context == nil {
			o.Context = context.Background()
		}
		if v, ok := o.Context.Value(codecsKey{}).(map[string]encoding.Codec); ok && v != nil {
			codecs = v
		}
		codecs[contentType] = c
		o.Context = context.WithValue(o.Context, codecsKey{}, codecs)
	}
}

// AuthTLS should be used to setup a secure authentication using TLS
func AuthTLS(t *tls.Config) server.Option {
	return setServerOption(tlsAuth{}, t)
}

// MaxConn specifies maximum number of max simultaneous connections to server
func MaxConn(n int) server.Option {
	return setServerOption(maxConnKey{}, n)
}

// Listener specifies the net.Listener to use instead of the default
func Listener(l net.Listener) server.Option {
	return setServerOption(netListener{}, l)
}

// Options to be used to configure gRPC options
func Options(opts ...grpc.ServerOption) server.Option {
	return setServerOption(grpcOptions{}, opts)
}

// GRPCWebOptions to be used to start a gRPC Web server
func GRPCWebOptions(opts ...grpcweb.Option) server.Option {
	return setServerOption(grpcWebOptions{}, opts)
}

// GRPCWebPort to be used to start a gRPC Web server
func GRPCWebPort(addr string) server.Option {
	return setServerOption(grpcWebPort{}, addr)
}

//
// Deprecated: use MaxRecvMsgSize or MaxSendMsgSize instead
// MaxMsgSize set the maximum message in bytes the server can receive and
// send.  Default maximum message size is 4 MB.
func MaxMsgSize(s int) server.Option {
	return setServerOption(maxMsgSizeKey{}, s)
}

//
// MaxRecvMsgSize set the maximum size of message that server can receive.
//
func MaxRecvMsgSize(s int) server.Option {
	return func(o *server.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, maxRecvMsgSizeKey{}, s)
	}
}

//
// MaxSendMsgSize set the maximum size of message that server can send.
//
func MaxSendMsgSize(s int) server.Option {
	return func(o *server.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, maxSendMsgSizeKey{}, s)
	}
}

func newOptions(opt ...server.Option) server.Options {
	opts := server.Options{
		Codecs:           make(map[string]codec.NewCodec),
		Metadata:         map[string]string{},
		Broker:           memory.NewBroker(),
		Registry:         mdns.NewRegistry(),
		Address:          server.DefaultAddress,
		Name:             server.DefaultName,
		Id:               server.DefaultId,
		Version:          server.DefaultVersion,
		RegisterInterval: server.DefaultRegisterInterval,
		RegisterTTL:      server.DefaultRegisterTTL,
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}
