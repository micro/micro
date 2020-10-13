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
// Original source: github.com/micro/go-micro/v3/network/transport/grpc/socket.go

package grpc

import (
	"github.com/micro/micro/v3/internal/network/transport"
	pb "github.com/micro/micro/v3/proto/transport"
	"google.golang.org/grpc"
)

type grpcTransportClient struct {
	conn   *grpc.ClientConn
	stream pb.Transport_StreamClient

	local  string
	remote string
}

type grpcTransportSocket struct {
	stream pb.Transport_StreamServer
	local  string
	remote string
}

func (g *grpcTransportClient) Local() string {
	return g.local
}

func (g *grpcTransportClient) Remote() string {
	return g.remote
}

func (g *grpcTransportClient) Recv(m *transport.Message) error {
	if m == nil {
		return nil
	}

	msg, err := g.stream.Recv()
	if err != nil {
		return err
	}

	m.Header = msg.Header
	m.Body = msg.Body
	return nil
}

func (g *grpcTransportClient) Send(m *transport.Message) error {
	if m == nil {
		return nil
	}

	return g.stream.Send(&pb.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (g *grpcTransportClient) Close() error {
	return g.conn.Close()
}

func (g *grpcTransportSocket) Local() string {
	return g.local
}

func (g *grpcTransportSocket) Remote() string {
	return g.remote
}

func (g *grpcTransportSocket) Recv(m *transport.Message) error {
	if m == nil {
		return nil
	}

	msg, err := g.stream.Recv()
	if err != nil {
		return err
	}

	m.Header = msg.Header
	m.Body = msg.Body
	return nil
}

func (g *grpcTransportSocket) Send(m *transport.Message) error {
	if m == nil {
		return nil
	}

	return g.stream.Send(&pb.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (g *grpcTransportSocket) Close() error {
	return nil
}
