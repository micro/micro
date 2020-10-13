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
// Original source: github.com/micro/go-micro/v3/codec/proto/marshaler.go

package proto

import (
	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/micro/micro/v3/internal/codec"
	"github.com/oxtoacart/bpool"
)

// create buffer pool with 16 instances each preallocated with 256 bytes
var bufferPool = bpool.NewSizedBufferPool(16, 256)

type Marshaler struct{}

func (Marshaler) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, codec.ErrInvalidMessage
	}

	// looks not good, but allows to reuse underlining bytes
	buf := bufferPool.Get()
	pbuf := proto.NewBuffer(buf.Bytes())
	defer func() {
		bufferPool.Put(bytes.NewBuffer(pbuf.Bytes()))
	}()

	if err := pbuf.Marshal(pb); err != nil {
		return nil, err
	}

	return pbuf.Bytes(), nil
}

func (Marshaler) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return codec.ErrInvalidMessage
	}

	return proto.Unmarshal(data, pb)
}

func (Marshaler) String() string {
	return "proto"
}
