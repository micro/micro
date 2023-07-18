// Package model for managing data access
package model

import (
	"errors"
)

var (
	ErrorNilInterface         = errors.New("interface is nil")
	ErrorNotFound             = errors.New("not found")
	ErrorMultipleRecordsFound = errors.New("multiple records found")
)

var (
	DefaultModel Model
)

// Model represents a place where data can be saved to and
// queried from.
type Model interface {
	// Register a new model eg. User struct, Order struct
	Register(v interface{}) error
	// Create a new object. (Maintains indexes set up)
	Create(v interface{}) error
	// Update will take an existing object and update it.
	Update(v interface{}) error
	// Read a result by id e.g &User{ID: 1}
	Read(v interface{}) error
	// Deletes a record
	Delete(v interface{}) error
	// Query with a where clause
	Query(res interface{}, where ...interface{}) error
}

type Options struct {
	// Database sets the default database
	Database string
	// Table sets the default table
	Table string
	// Address to scope to
	Address string
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

// WithAddress sets the database address
func WithAddress(ns string) Option {
	return func(o *Options) {
		o.Address = ns
	}
}
