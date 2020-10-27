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
// Original source: github.com/micro/go-micro/v3/codec/proto/message.go

package proto

type Message struct {
	Data []byte
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return m.Data, nil
}

func (m *Message) UnmarshalJSON(data []byte) error {
	m.Data = data
	return nil
}

func (m *Message) ProtoMessage() {}

func (m *Message) Reset() {
	*m = Message{}
}

func (m *Message) String() string {
	return string(m.Data)
}

func (m *Message) Marshal() ([]byte, error) {
	return m.Data, nil
}

func (m *Message) Unmarshal(data []byte) error {
	m.Data = data
	return nil
}

func NewMessage(data []byte) *Message {
	return &Message{data}
}
