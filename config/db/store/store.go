package storeDB

import (
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/micro/v2/config/db"
)

type storeDB struct {
	opts db.Options
	st   store.Store
}

func init() {
	// new storeDB
	st := new(storeDB)
	// set default store
	st.st = store.DefaultStore
	// register
	db.Register(st)
}

func (m *storeDB) Init(opts db.Options) error {
	m.opts = opts
	// set the store we use
	if m.opts.Store != nil {
		m.st = m.opts.Store
	}
	return nil
}

func (m *storeDB) Create(record *store.Record) error {
	return m.st.Write(record)
}

func (m *storeDB) Read(key string) (*store.Record, error) {
	s, err := m.st.Read(key)
	if err != nil {
		return nil, err
	}

	return s[0], nil
}

func (m *storeDB) Update(record *store.Record) error {
	return m.st.Write(record)
}

func (m *storeDB) Delete(key string) error {
	return m.st.Delete(key)
}

func (m storeDB) List(opts ...db.ListOption) ([]*store.Record, error) {
	return m.st.List()
}

func (m *storeDB) String() string {
	return "store"
}
