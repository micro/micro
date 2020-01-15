package mysql

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/store"
	"github.com/micro/micro/config/db"
)

var (
	st         = map[string]*sql.Stmt{}
	defaultUrl = "root:123@(127.0.0.1:3306)/config?charset=utf8&parseTime=true"
)

type mysql struct {
	db *sql.DB
	st store.Store
}

func init() {
	db.Register(new(mysql))
}

func (m *mysql) Init(opts db.Options) error {
	var d *sql.DB
	var err error

	if opts.Url != "" {
		defaultUrl = opts.Url
	}

	parts := strings.Split(defaultUrl, "/")
	if len(parts) != 2 {
		return errors.New("Invalid database url ")
	}

	if len(parts[1]) == 0 {
		return errors.New("Invalid database name ")
	}

	var paramParts []string
	if strings.Contains(defaultUrl, "?") {
		paramParts = strings.Split(parts[1], "?")
		parts[1] = paramParts[0]
		paramParts = paramParts[1:]
	}

	url := parts[0]
	database := "`" + parts[1] + "`"

	if d, err = sql.Open("mysql", url+"/"); err != nil {
		return err
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		return err
	}
	d.Close()

	if d, err = sql.Open("mysql", defaultUrl); err != nil {
		return err
	}
	if _, err = d.Exec(changeSchema); err != nil {
		return err
	}

	m.db = d

	return nil
}

func (m *mysql) Create(record *store.Record) error {
	return m.st.Write(record)
}

func (m *mysql) Read(key string) (*store.Record, error) {
	s, err := m.st.Read(key)
	if err != nil {
		return nil, err
	}

	return s[0], nil
}

func (m *mysql) Delete(key string) error {
	return m.st.Delete(key)
}

func (m *mysql) Update(record *store.Record) error {
	return m.st.Write(record)
}

func (m *mysql) List(opts ...db.ListOption) ([]*store.Record, error) {
	// opts is just params holder
	return m.st.List()
}

func (m *mysql) String() string {
	return "mysql"
}
