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
// Original source: github.com/micro/go-micro/v3/client/mucp/mucp_response.go

package mucp

import (
	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/internal/network/transport"
)

type rpcResponse struct {
	header map[string]string
	body   []byte
	socket transport.Socket
	codec  codec.Codec
}

func (r *rpcResponse) Codec() codec.Reader {
	return r.codec
}

func (r *rpcResponse) Header() map[string]string {
	return r.header
}

func (r *rpcResponse) Read() ([]byte, error) {
	var msg transport.Message

	if err := r.socket.Recv(&msg); err != nil {
		return nil, err
	}

	// set internals
	r.header = msg.Header
	r.body = msg.Body

	return msg.Body, nil
}
