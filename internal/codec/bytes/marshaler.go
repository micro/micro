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
// Original source: github.com/micro/go-micro/v3/codec/bytes/marshaler.go

package bytes

import (
	"github.com/micro/micro/v3/internal/codec"
)

type Marshaler struct{}

type Message struct {
	Header map[string]string
	Body   []byte
}

func (n Marshaler) Marshal(v interface{}) ([]byte, error) {
	switch ve := v.(type) {
	case *[]byte:
		return *ve, nil
	case []byte:
		return ve, nil
	case *Message:
		return ve.Body, nil
	}
	return nil, codec.ErrInvalidMessage
}

func (n Marshaler) Unmarshal(d []byte, v interface{}) error {
	switch ve := v.(type) {
	case *[]byte:
		*ve = d
	case *Message:
		ve.Body = d
	}
	return codec.ErrInvalidMessage
}

func (n Marshaler) String() string {
	return "bytes"
}
