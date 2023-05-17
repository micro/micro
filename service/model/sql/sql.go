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

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/util/user"
)

// The sql DB
type DB struct {
	*gorm.DB
}

// NewDB provides a new database connection. If [name].db.address is found
// in the config then it's used as the address, otherwise we use sqlite.
func NewDB(name string) (*DB, error) {
	dbAddress := "sqlite://" + name + ".db"

	// Connect to the database
	cfg, err := config.Get(name + ".db.address")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)

	var db *gorm.DB

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
