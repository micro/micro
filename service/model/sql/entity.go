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
// Original source: github.com/micro/go-micro/v3/model/sql/entity.go

package sql

import (
	"github.com/google/uuid"
	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/service/model"
)

type sqlEntity struct {
	id         string
	name       string
	value      interface{}
	codec      codec.Marshaler
	attributes map[string]interface{}
}

func (m *sqlEntity) Attributes() map[string]interface{} {
	return m.attributes
}

func (m *sqlEntity) Id() string {
	return m.id
}

func (m *sqlEntity) Name() string {
	return m.name
}

func (m *sqlEntity) Value() interface{} {
	return m.value
}

func newEntity(name string, value interface{}, codec codec.Marshaler) model.Entity {
	return &sqlEntity{
		id:         uuid.New().String(),
		name:       name,
		value:      value,
		codec:      codec,
		attributes: make(map[string]interface{}),
	}
}
