package cockroach

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	proto "github.com/micro/go-micro/config/source/mucp/proto"
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

func (m *cockroach) Create(change *proto.Change) error {
	rd := &store.Record{
		Key:   change.Key,
		Value: change.ChangeSet.Data,
	}
	m.st.Write()
	return nil
}

func (m *cockroach) Read(id string) (*proto.Change, error) {
	if len(id) == 0 {
		return nil, errors.New("Invalid trace id")
	}

	return nil, nil
}

func (m *cockroach) Delete(change *proto.Change) error {

	return nil
}

func (m *cockroach) Update(change *proto.Change) error {

	return nil
}

func (m *cockroach) List(opts db.ListOptions) ([]*proto.Change, error) {
	return nil, nil
}

func (m *cockroach) String() string {
	return "cockroach"
}
