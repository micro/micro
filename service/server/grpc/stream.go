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
// Original source: github.com/micro/go-micro/v3/server/grpc/stream.go

package grpc

import (
	"context"

	"github.com/micro/micro/v3/service/server"
	"google.golang.org/grpc"
)

// rpcStream implements a server side Stream.
type rpcStream struct {
	// embed the grpc stream so we can access it
	grpc.ServerStream

	request server.Request
}

func (r *rpcStream) Close() error {
	return nil
}

func (r *rpcStream) Error() error {
	return nil
}

func (r *rpcStream) Request() server.Request {
	return r.request
}

func (r *rpcStream) Context() context.Context {
	return r.ServerStream.Context()
}

func (r *rpcStream) Send(m interface{}) error {
	return r.ServerStream.SendMsg(m)
}

func (r *rpcStream) Recv(m interface{}) error {
	return r.ServerStream.RecvMsg(m)
}
