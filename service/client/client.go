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
// Original source: github.com/micro/go-micro/v3/client/client.go

// Package client is an interface for an RPC client
package client

import (
	"context"
	"time"

	"github.com/micro/micro/v5/util/codec"
)

var (
	// DefaultClient implementation
	DefaultClient Client
)

// Client is the interface used to make requests to services.
// It supports Request/Response via Transport and Publishing via the Broker.
// It also supports bidirectional streaming of requests.
type Client interface {
	Init(...Option) error
	Options() Options
	NewMessage(topic string, msg interface{}, opts ...MessageOption) Message
	NewRequest(service, endpoint string, req interface{}, reqOpts ...RequestOption) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
	Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
	Publish(ctx context.Context, msg Message, opts ...PublishOption) error
	String() string
}

// Message is the interface for publishing asynchronously
type Message interface {
	Topic() string
	Payload() interface{}
	ContentType() string
}

// Request is the interface for a synchronous request used by Call or Stream
type Request interface {
	// The service to call
	Service() string
	// The action to take
	Method() string
	// The endpoint to invoke
	Endpoint() string
	// The content type
	ContentType() string
	// The unencoded request body
	Body() interface{}
	// Write to the encoded request writer. This is nil before a call is made
	Codec() codec.Writer
	// indicates whether the request will be a streaming one rather than unary
	Stream() bool
}

// Response is the response received from a service
type Response interface {
	// Read the response
	Codec() codec.Reader
	// read the header
	Header() map[string]string
	// Read the undecoded response
	Read() ([]byte, error)
}

// Stream is the interface for a bidirectional synchronous stream
type Stream interface {
	// Context for the stream
	Context() context.Context
	// The request made
	Request() Request
	// The response read
	Response() Response
	// Send will encode and send a request
	Send(interface{}) error
	// Recv will decode and read a response
	Recv(interface{}) error
	// Error returns the stream error
	Error() error
	// Close closes the stream
	Close() error
}

// Option used by the Client
type Option func(*Options)

// CallOption used by Call or Stream
type CallOption func(*CallOptions)

// PublishOption used by Publish
type PublishOption func(*PublishOptions)

// MessageOption used by NewMessage
type MessageOption func(*MessageOptions)

// RequestOption used by NewRequest
type RequestOption func(*RequestOptions)

var (
	// DefaultBackoff is the default backoff function for retries
	DefaultBackoff = exponentialBackoff
	// DefaultRetry is the default check-for-retry function for retries
	DefaultRetry = RetryOnConnectFailure
	// DefaultRetries is the default number of times a request is tried
	DefaultRetries = 1
	// DefaultRequestTimeout is the default request timeout
	DefaultRequestTimeout = time.Minute
	// DefaultPoolSize sets the connection pool size
	DefaultPoolSize = 100
	// DefaultPoolTTL sets the connection pool ttl
	DefaultPoolTTL = time.Minute
)

// NewMessage returns a message which can be published
func NewMessage(topic string, msg interface{}) Message {
	return DefaultClient.NewMessage(topic, msg)
}

// NewRequest returns a request can which be executed using Call or Stream
func NewRequest(service, endpoint string, req interface{}, opts ...RequestOption) Request {
	return DefaultClient.NewRequest(service, endpoint, req, opts...)
}

// Call performs a request
func Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {
	return DefaultClient.Call(ctx, req, rsp, opts...)
}

// Publish a message
func Publish(ctx context.Context, msg Message) error {
	return DefaultClient.Publish(ctx, msg)
}
