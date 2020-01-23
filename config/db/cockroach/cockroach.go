package cockroach

import (
	"database/sql"
	"net/url"
	"strings"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/store"
	roachStore "github.com/micro/go-micro/store/cockroach"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/config/db"
)

var (
	defaultUrl = "postgres://root:@127.0.0.1:26257/micro?search_path=config"
)

type cockroach struct {
	db *sql.DB
	st store.Store
}

func init() {
	db.Register(new(cockroach))
}

func (m *cockroach) Init(opts db.Options) error {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	var d *sql.DB

	if opts.Url != "" {
		defaultUrl = opts.Url
	}

	u, _ := url.Parse(defaultUrl)
	parts := strings.Split(u.Path, "/")
	if len(parts) != 2 {
		err = errors.InternalServerError("go.micro.config.Init", "Invalid database url for cockroach: %s", u.Path)
		return err
	}

	schema := u.Query().Get("search_path")
	if len(schema) == 0 {
		err = errors.InternalServerError("go.micro.config.Init", "needs a schema with search_path")
		return err
	}

	m.db = d
	m.st = roachStore.NewStore(
		store.Nodes(opts.Url),
		store.Namespace(schema))

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

func (m *cockroach) List(opts ...db.ListOption) ([]*store.Record, error) {
	return m.st.List()
}

func (m *cockroach) String() string {
	return "cockroach"
}
