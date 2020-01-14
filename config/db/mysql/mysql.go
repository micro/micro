package mysql

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	proto "github.com/micro/go-micro/config/source/mucp/proto"
	"github.com/micro/micro/config/db"
)

var (
	st = map[string]*sql.Stmt{}
)

type mysql struct {
	db *sql.DB
}

func init() {
	db.Register(new(mysql))
}

func (m *mysql) Init(opts db.Options) error {
	var d *sql.DB
	var err error

	parts := strings.Split(opts.Url, "/")
	if len(parts) != 2 {
		return errors.New("Invalid database url ")
	}

	if len(parts[1]) == 0 {
		return errors.New("Invalid database name ")
	}

	var paramParts []string
	if strings.Contains(opts.Url, "?") {
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

	if d, err = sql.Open("mysql", opts.Url); err != nil {
		return err
	}
	if _, err = d.Exec(changeSchema); err != nil {
		return err
	}

	m.db = d

	return nil
}

func (m *mysql) Create(change *proto.Change) error {

	return nil
}

func (m *mysql) Read(id string) (*proto.Change, error) {
	if len(id) == 0 {
		return nil, errors.New("Invalid trace id")
	}

	return nil, nil
}

func (m *mysql) Delete(change *proto.Change) error {

	return nil
}

func (m *mysql) Update(change *proto.Change) error {

	return nil
}

func (m *mysql) List(opts db.ListOptions) ([]*proto.Change, error) {
	return nil, nil
}

func (m *mysql) String() string {
	return "mysql"
}
