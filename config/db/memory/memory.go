package memory

import (
	"github.com/micro/go-micro/store"
	mStore "github.com/micro/go-micro/store/memory"
	"github.com/micro/micro/config/db"
)

type memory struct {
	st store.Store
}

func init() {
	db.Register(new(memory))
}

func (m *memory) Init(opts db.Options) error {
	m.st = mStore.NewStore()
	return nil
}

func (m *memory) Create(record *store.Record) error {
	return m.st.Write(record)
}

func (m *memory) Read(key string) (*store.Record, error) {
	s, err := m.st.Read(key)
	if err != nil {
		return nil, err
	}

	return s[0], nil
}

func (m *memory) Update(record *store.Record) error {
	return m.st.Write(record)
}

func (m *memory) Delete(key string) error {
	return m.st.Delete(key)
}

func (m memory) List(opts ...db.ListOption) ([]*store.Record, error) {
	return m.st.List()
}

func (m *memory) String() string {
	return "memory"
}
