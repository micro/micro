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
// Original source: github.com/micro/go-micro/v3/model/mud/entity.go

package mud

import (
	"github.com/google/uuid"
	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/service/model"
)

type mudEntity struct {
	id         string
	name       string
	value      interface{}
	codec      codec.Marshaler
	attributes map[string]interface{}
}

func (m *mudEntity) Attributes() map[string]interface{} {
	return m.attributes
}

func (m *mudEntity) Id() string {
	return m.id
}

func (m *mudEntity) Name() string {
	return m.name
}

func (m *mudEntity) Value() interface{} {
	return m.value
}

func (m *mudEntity) Read(v interface{}) error {
	switch m.value.(type) {
	case []byte:
		b := m.value.([]byte)
		return m.codec.Unmarshal(b, v)
	default:
		v = m.value
	}
	return nil
}

func newEntity(name string, value interface{}, codec codec.Marshaler) model.Entity {
	return &mudEntity{
		id:         uuid.New().String(),
		name:       name,
		value:      value,
		codec:      codec,
		attributes: make(map[string]interface{}),
	}
}
