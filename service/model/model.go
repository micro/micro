// Package model implements convenience methods for
// managing indexes on top of the Store.
// See this doc for the general idea https://github.com/m3o/dev/blob/feature/storeindex/design/auto-indexes.md
// Prior art/Inspirations from github.com/gocassa/gocassa, which
// is a similar package on top an other KV store (Cassandra/gocql)
package model

import (
	"context"
	"errors"

	"github.com/micro/micro/v3/service/store"
)

var (
	ErrorNilInterface         = errors.New("interface is nil")
	ErrorNotFound             = errors.New("not found")
	ErrorMultipleRecordsFound = errors.New("multiple records found")
)

type OrderType string

const (
	OrderTypeUnordered = OrderType("unordered")
	OrderTypeAsc       = OrderType("ascending")
	OrderTypeDesc      = OrderType("descending")
)

const (
	queryTypeEq  = "eq"
	indexTypeEq  = "eq"
	queryTypeAll = "all"
	indexTypeAll = "all"
)

var (
	// DefaultKey is the default field for indexing
	DefaultKey = "ID"

	// DefaultIndex is the ID index
	DefaultIndex = newIndex("ID")

	// DefaultModel is the default model
	DefaultModel = NewModel()
)

// Model represents a place where data can be saved to and
// queried from.
type Model interface {
	// Context sets the context for the model returning a new copy
	Context(ctx context.Context) Model
	// Register a new model eg. User struct, Order struct
	Register(v interface{}) error
	// Create a new object. (Maintains indexes set up)
	Create(v interface{}) error
	// Update will take an existing object and update it.
	// TODO: Make use of "sync" interface to lock, read, write, unlock
	Update(v interface{}) error
	// Read accepts a pointer to a value and expects to fine one or more
	// elements. Read throws an error if a value is not found or we can't
	// find a matching index for a slice based query.
	Read(query Query, resultPointer interface{}) error
	// Deletes a record. Delete only support Equals("id", value) for now.
	// @todo Delete only supports string keys for now.
	Delete(query Query) error
}

type Options struct {
	// Database sets the default database
	Database string
	// Table sets the default table
	Table string
	// Enable debug logging
	Debug bool
	// The indexes to use for queries
	Indexes []Index
	// Namespace to scope to
	Namespace string
	// Store is the storage engine
	Store store.Store
	// Context is the context for all model queries
	Context context.Context
	// Key is the fiel name of the primary key
	Key string
}

type Option func(*Options)

// WithDatabase sets the default database for queries
func WithDatabase(db string) Option {
	return func(o *Options) {
		o.Database = db
	}
}

// WithTable sets the default table for queries
func WithTable(t string) Option {
	return func(o *Options) {
		o.Table = t
	}
}

// WithContext sets the context for all queries
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// WithIndexes creates an option with the given indexes
func WithIndexes(idx ...Index) Option {
	return func(o *Options) {
		o.Indexes = idx
	}
}

// WithStore create an option for setting the store
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithDebug enables debug logging
func WithDebug(d bool) Option {
	return func(o *Options) {
		o.Debug = d
	}
}

// WithNamespace sets the namespace to scope to
func WithNamespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}

// WithKey sets the Key
func WithKey(idField string) Option {
	return func(o *Options) {
		o.Key = idField
	}
}
