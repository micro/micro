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
// Original source: github.com/micro/go-micro/v3/client/grpc/message.go

package grpc

import (
	"github.com/micro/micro/v3/service/client"
)

type grpcEvent struct {
	topic       string
	contentType string
	payload     interface{}
}

func newGRPCEvent(topic string, payload interface{}, contentType string, opts ...client.MessageOption) client.Message {
	var options client.MessageOptions
	for _, o := range opts {
		o(&options)
	}

	if len(options.ContentType) > 0 {
		contentType = options.ContentType
	}

	return &grpcEvent{
		payload:     payload,
		topic:       topic,
		contentType: contentType,
	}
}

func (g *grpcEvent) ContentType() string {
	return g.contentType
}

func (g *grpcEvent) Topic() string {
	return g.topic
}

func (g *grpcEvent) Payload() interface{} {
	return g.payload
}
