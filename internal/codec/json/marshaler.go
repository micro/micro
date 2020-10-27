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
// Original source: github.com/micro/go-micro/v3/codec/json/marshaler.go

package json

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/oxtoacart/bpool"
)

var jsonpbMarshaler = &jsonpb.Marshaler{}

// create buffer pool with 16 instances each preallocated with 256 bytes
var bufferPool = bpool.NewSizedBufferPool(16, 256)

type Marshaler struct{}

func (j Marshaler) Marshal(v interface{}) ([]byte, error) {
	if pb, ok := v.(proto.Message); ok {
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)
		if err := jsonpbMarshaler.Marshal(buf, pb); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return json.Marshal(v)
}

func (j Marshaler) Unmarshal(d []byte, v interface{}) error {
	if pb, ok := v.(proto.Message); ok {
		return jsonpb.Unmarshal(bytes.NewReader(d), pb)
	}
	return json.Unmarshal(d, v)
}

func (j Marshaler) String() string {
	return "json"
}
