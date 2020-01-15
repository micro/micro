package cockroach

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/store"
	pgStore "github.com/micro/go-micro/store/cockroach"
	"github.com/micro/micro/config/db"
)

var (
	defaultUrl = "postgres://postgres:@localhost:5432/postgres"
)

type cockroach struct {
	db *sql.DB
	st store.Store
}

func init() {
	db.Register(new(cockroach))
}

func (m *cockroach) Init(opts db.Options) error {
	var d *sql.DB
	var err error

	if opts.Url == "" {
		defaultUrl = opts.Url
	}

	if d, err = sql.Open("postgres", defaultUrl); err != nil {
		return err
	}

	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + opts.DBName); err != nil {
		return err
	}
	d.Close()

	if d, err = sql.Open("cockroach", opts.Url); err != nil {
		return err
	}
	if _, err = d.Exec(changeSchema); err != nil {
		return err
	}

	m.db = d
	m.st = pgStore.NewStore(store.Nodes(opts.Url))

	return nil
}

func (m *cockroach) Create(record *store.Record) error {
	return m.st.Write(record)
}

func (m *cockroach) Read(key string) (*store.Record, error) {
	s, err := m.st.Read(key)
	if err != nil {
		return nil, err
	}

	return s[0], nil
}

func (m *cockroach) Delete(key string) error {
	return m.st.Delete(key)
}

func (m *cockroach) Update(record *store.Record) error {
	return m.st.Write(record)
}

func (m *cockroach) List(opts db.ListOptions) ([]*store.Record, error) {
	return m.st.List()
}

func (m *cockroach) String() string {
	return "cockroach"
}
