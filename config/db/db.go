package db

import (
	"errors"
	"sync"

	"github.com/micro/go-micro/util/log"
	proto "github.com/micro/micro/config/proto/config"
)

type DB interface {
	Init(Options) error
	Create(*proto.Change) error
	Read(id string) (*proto.Change, error)
	Update(*proto.Change) error
	Delete(*proto.Change) error
	Search(id, author string, limit, offset int64) ([]*proto.Change, error)
	AuditLog(from, to int64, limit, offset int64, reverse bool) ([]*proto.ChangeLog, error)
	String() string
}

var (
	db          DB
	dbMap       = map[string]DB{}
	mux         sync.Mutex
	ErrNotFound = errors.New("not found")
)

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

func Create(ch *proto.Change) error {
	return db.Create(ch)
}

func Read(id string) (*proto.Change, error) {
	return db.Read(id)
}

func Update(ch *proto.Change) error {
	return db.Update(ch)
}

func Delete(ch *proto.Change) error {
	return db.Delete(ch)
}

func Search(id, author string, limit, offset int64) ([]*proto.Change, error) {
	return db.Search(id, author, limit, offset)
}

func AuditLog(from, to, limit, offset int64, reverse bool) ([]*proto.ChangeLog, error) {
	return db.AuditLog(from, to, limit, offset, reverse)
}
