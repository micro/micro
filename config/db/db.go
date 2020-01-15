package db

import (
	"errors"
	"sync"

	"github.com/micro/go-micro/store"
	"github.com/micro/go-micro/util/log"
)

var (
	db          DB
	dbMap       = map[string]DB{}
	mux         sync.Mutex
	ErrNotFound = errors.New("not found")
)

type DB interface {
	Init(Options) error
	Create(*store.Record) error
	Read(key string) (*store.Record, error)
	Update(*store.Record) error
	Delete(key string) error
	List(opts ...ListOption) ([]*store.Record, error)
	String() string
}

func Register(backend DB) {
	mux.Lock()
	defer mux.Unlock()

	if dbMap[backend.String()] != nil {
		log.Fatalf("db is repeated: %s", backend.String())
	}

	dbMap[backend.String()] = backend
	log.Logf("Register config db: %s", backend.String())
}

func Init(opts ...Option) error {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	db = dbMap[options.DBName]
	return db.Init(options)
}

func Create(ch *store.Record) error {
	return db.Create(ch)
}

func Read(key string) (*store.Record, error) {
	return db.Read(key)
}

func Update(ch *store.Record) error {
	return db.Update(ch)
}

func Delete(key string) error {
	return db.Delete(key)
}

func List(opts ...ListOption) ([]*store.Record, error) {
	return db.List(opts...)
}
