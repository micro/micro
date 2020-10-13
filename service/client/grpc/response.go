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
// Original source: github.com/micro/go-micro/v3/client/grpc/response.go

package grpc

import (
	"strings"

	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/internal/codec/bytes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

type response struct {
	conn   *poolConn
	stream grpc.ClientStream
	codec  encoding.Codec
	gcodec codec.Codec
}

// Read the response
func (r *response) Codec() codec.Reader {
	return r.gcodec
}

// read the header
func (r *response) Header() map[string]string {
	md, err := r.stream.Header()
	if err != nil {
		return map[string]string{}
	}
	hdr := make(map[string]string, len(md))
	for k, v := range md {
		hdr[k] = strings.Join(v, ",")
	}
	return hdr
}

// Read the undecoded response
func (r *response) Read() ([]byte, error) {
	f := &bytes.Frame{}
	if err := r.gcodec.ReadBody(f); err != nil {
		return nil, err
	}
	return f.Data, nil
}
