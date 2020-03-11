package db

import (
	"errors"
	"sync"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
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
}

func Init(opts ...Option) error {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	db = dbMap[options.Database]

	// initialise db config
	log.Infof("Init config options: %+v", options)

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
