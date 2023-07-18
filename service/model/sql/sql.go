// Package sql is for a sql database model
package sql

import (
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"micro.dev/v4/service/model"
	"micro.dev/v4/util/user"
)

// The sql DB
type DB struct {
	*gorm.DB
}

func (d *DB) Register(v interface{}) error {
	return d.DB.AutoMigrate(v)
}

func (d *DB) Create(v interface{}) error {
	return d.DB.Create(v).Error
}

func (d *DB) Update(v interface{}) error {
	return d.DB.Save(v).Error
}

func (d *DB) Delete(v interface{}) error {
	return d.DB.Delete(v).Error
}

func (d *DB) Read(v interface{}) error {
	return d.DB.First(v).Error
}

func (d *DB) Query(res interface{}, where ...interface{}) error {
	return d.DB.Find(res, where...).Error
}

func NewModel(opts ...model.Option) model.Model {
	var options model.Options
	for _, o := range opts {
		o(&options)
	}

	if len(options.Database) == 0 {
		options.Database = "micro"
	}

	// create a new database handle
	db, _ := NewDB(options.Database, options.Address)

	return db
}

// NewDB provides a new database connection. If [name].db.address is found
// in the config then it's used as the address, otherwise we use sqlite.
func NewDB(name, addr string) (*DB, error) {
	var db *gorm.DB
	var err error

	if strings.HasPrefix(addr, "postgres") {
		db, err = gorm.Open(postgres.Open(addr), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: name + "_",
			},
		})
	} else {
		path := filepath.Join(user.Dir, "service", name)
		os.MkdirAll(path, 0755)
		file := filepath.Join(path, "db.sqlite")
		db, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	}

	return &DB{db}, err
}
