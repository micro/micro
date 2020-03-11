package db

import (
	"github.com/micro/go-micro/v2/store"
)

type Options struct {
	Url      string
	Database string
	Table    string
	// Store
	Store store.Store
}

type Option func(*Options)

func WithStore(st store.Store) Option {
	return func(o *Options) {
		o.Store = st
	}
}

// WithDatabase sets which database store to use e.g memory, etc, cockroach, store
func WithDatabase(name string) Option {
	return func(options *Options) {
		options.Database = name
	}
}

// WithTable set the table to store data, if supported.
func WithTable(table string) Option {
	return func(options *Options) {
		options.Table = table
	}
}

func WithUrl(url string) Option {
	return func(options *Options) {
		options.Url = url
	}
}

type ListOptions struct{}

type ListOption func(*Options)
