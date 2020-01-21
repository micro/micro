package etcd

import (
	"github.com/micro/go-micro/store"
	eStore "github.com/micro/go-micro/store/etcd"
	"github.com/micro/micro/config/db"
	"strings"
)

var (
	defaultUrl = "http://127.0.0.1:2379"
)

type etcd struct {
	st store.Store
}

func init() {
	db.Register(new(etcd))
}

func (m *etcd) Init(opts db.Options) error {
	if opts.Url != "" {
		defaultUrl = opts.Url
	}

	m.st = eStore.NewStore(store.Nodes(strings.Split(defaultUrl, ",")...))
	return nil
}

func (m *etcd) Create(record *store.Record) error {
	return m.st.Write(record)
}

func (m *etcd) Read(key string) (*store.Record, error) {
	s, err := m.st.Read(key)
	if err != nil {
		return nil, err
	}

	return s[0], nil
}

func (m *etcd) Update(record *store.Record) error {
	return m.st.Write(record)
}

func (m *etcd) Delete(key string) error {
	return m.st.Delete(key)
}

func (m etcd) List(opts ...db.ListOption) ([]*store.Record, error) {
	return m.st.List()
}

func (m *etcd) String() string {
	return "etcd"
}
