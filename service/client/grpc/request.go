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
// Original source: github.com/micro/go-micro/v3/client/grpc/request.go

package grpc

import (
	"fmt"
	"strings"

	"github.com/micro/micro/v5/service/client"
	"github.com/micro/micro/v5/util/codec"
)

type grpcRequest struct {
	service     string
	method      string
	contentType string
	request     interface{}
	opts        client.RequestOptions
	codec       codec.Codec
}

// service Struct.Method /service.Struct/Method
func methodToGRPC(service, method string) string {
	// no method or already grpc method
	if len(method) == 0 || method[0] == '/' {
		return method
	}

	// assume method is Foo.Bar
	mParts := strings.Split(method, ".")
	if len(mParts) != 2 {
		return method
	}

	if len(service) == 0 {
		return fmt.Sprintf("/%s/%s", mParts[0], mParts[1])
	}

	// return /pkg.Foo/Bar
	return fmt.Sprintf("/%s.%s/%s", service, mParts[0], mParts[1])
}

func newGRPCRequest(service, method string, request interface{}, contentType string, reqOpts ...client.RequestOption) client.Request {
	var opts client.RequestOptions
	for _, o := range reqOpts {
		o(&opts)
	}

	// set the content-type specified
	if len(opts.ContentType) > 0 {
		contentType = opts.ContentType
	}

	return &grpcRequest{
		service:     service,
		method:      method,
		request:     request,
		contentType: contentType,
		opts:        opts,
	}
}

func (g *grpcRequest) ContentType() string {
	return g.contentType
}

func (g *grpcRequest) Service() string {
	return g.service
}

func (g *grpcRequest) Method() string {
	return g.method
}

func (g *grpcRequest) Endpoint() string {
	return g.method
}

func (g *grpcRequest) Codec() codec.Writer {
	return g.codec
}

func (g *grpcRequest) Body() interface{} {
	return g.request
}

func (g *grpcRequest) Stream() bool {
	return g.opts.Stream
}
