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
// Original source: github.com/micro/go-micro/v3/server/grpc/request.go

package grpc

import (
	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/internal/codec/bytes"
)

type rpcRequest struct {
	service     string
	method      string
	contentType string
	codec       codec.Codec
	header      map[string]string
	body        []byte
	stream      bool
	payload     interface{}
}

type rpcMessage struct {
	topic       string
	contentType string
	payload     interface{}
	header      map[string]string
	body        []byte
	codec       codec.Codec
}

func (r *rpcRequest) ContentType() string {
	return r.contentType
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Endpoint() string {
	return r.method
}

func (r *rpcRequest) Codec() codec.Reader {
	return r.codec
}

func (r *rpcRequest) Header() map[string]string {
	return r.header
}

func (r *rpcRequest) Read() ([]byte, error) {
	f := &bytes.Frame{}
	if err := r.codec.ReadBody(f); err != nil {
		return nil, err
	}
	return f.Data, nil
}

func (r *rpcRequest) Stream() bool {
	return r.stream
}

func (r *rpcRequest) Body() interface{} {
	return r.payload
}

func (r *rpcMessage) ContentType() string {
	return r.contentType
}

func (r *rpcMessage) Topic() string {
	return r.topic
}

func (r *rpcMessage) Payload() interface{} {
	return r.payload
}

func (r *rpcMessage) Header() map[string]string {
	return r.header
}

func (r *rpcMessage) Body() []byte {
	return r.body
}

func (r *rpcMessage) Codec() codec.Reader {
	return r.codec
}
