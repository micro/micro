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
// Original source: github.com/micro/go-micro/v3/model/sql/sql.go

// Package sql is the micro data model implementation
package sql

import (
	"github.com/micro/go-micro/v3/codec/json"
	"github.com/micro/go-micro/v3/model"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/go-micro/v3/store/memory"
	memsync "github.com/micro/go-micro/v3/sync/memory"
)

type sqlModel struct {
	options model.Options
}

func (m *sqlModel) Init(opts ...model.Option) error {
	for _, o := range opts {
		o(&m.options)
	}
	return nil
}

func (m *sqlModel) NewEntity(name string, value interface{}) model.Entity {
	// TODO: potentially pluralise name for tables
	return newEntity(name, value, m.options.Codec)
}

func (m *sqlModel) Create(e model.Entity) error {
	// lock on the name of entity
	if err := m.options.Sync.Lock(e.Name()); err != nil {
		return err
	}
	// TODO: deal with the error
	defer m.options.Sync.Unlock(e.Name())

	// TODO: potentially add encode to entity?
	v, err := m.options.Codec.Marshal(e.Value())
	if err != nil {
		return err
	}

	// TODO: include metadata and set database
	return m.options.Store.Write(&store.Record{
		Key:   e.Id(),
		Value: v,
	}, store.WriteTo(m.options.Database, e.Name()))
}

func (m *sqlModel) Read(opts ...model.ReadOption) ([]model.Entity, error) {
	var options model.ReadOptions
	for _, o := range opts {
		o(&options)
	}
	// TODO: implement the options that allow querying
	return nil, nil
}

func (m *sqlModel) Update(e model.Entity) error {
	// TODO: read out the record first, update the fields and store

	// lock on the name of entity
	if err := m.options.Sync.Lock(e.Name()); err != nil {
		return err
	}
	// TODO: deal with the error
	defer m.options.Sync.Unlock(e.Name())

	// TODO: potentially add encode to entity?
	v, err := m.options.Codec.Marshal(e.Value())
	if err != nil {
		return err
	}

	// TODO: include metadata and set database
	return m.options.Store.Write(&store.Record{
		Key:   e.Id(),
		Value: v,
	}, store.WriteTo(m.options.Database, e.Name()))
}

func (m *sqlModel) Delete(opts ...model.DeleteOption) error {
	var options model.DeleteOptions
	for _, o := range opts {
		o(&options)
	}
	// TODO: implement the options that allow deleting
	return nil
}

func (m *sqlModel) String() string {
	return "sql"
}

func NewModel(opts ...model.Option) model.Model {
	options := model.Options{
		Codec: new(json.Marshaler),
		Sync:  memsync.NewSync(),
		Store: memory.NewStore(),
	}

	return &sqlModel{
		options: options,
	}
}
