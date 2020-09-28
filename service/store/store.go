package store

import (
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/store/client"
)

var (
	// DefaultStore implementation
	DefaultStore store.Store = client.NewStore()
	// DefaultBlobStore implementation
	DefaultBlobStore store.BlobStore = client.NewBlobStore()
	// ErrNotFound is returned when a key doesn't exist
	ErrNotFound = store.ErrNotFound
	// ErrMissingKey is returned when a key wasn't provided
	ErrMissingKey = store.ErrMissingKey
)

type (
	// Record is an alias for store.Record
	Record = store.Record
)

// Read records
func Read(key string, opts ...Option) ([]*Record, error) {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	// convert the options
	var readOpts []store.ReadOption
	if len(options.Prefix) > 0 {
		key = options.Prefix
		readOpts = append(readOpts, store.ReadPrefix())
	}
	if options.Limit > 0 {
		readOpts = append(readOpts, store.ReadLimit(options.Limit))
	}
	if options.Offset > 0 {
		readOpts = append(readOpts, store.ReadOffset(options.Offset))
	}

	// execute the query
	return DefaultStore.Read(key, readOpts...)
}

// Write a record to the store
func Write(r *Record) error {
	return DefaultStore.Write(r)
}

// Delete removes the record with the corresponding key from the store.
func Delete(key string) error {
	return DefaultStore.Delete(key)
}

// List returns any keys that match, or an empty list with no error if none matched.
func List(opts ...Option) ([]string, error) {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	// convert the options
	var listOpts []store.ListOption
	if len(options.Prefix) > 0 {
		listOpts = append(listOpts, store.ListPrefix(options.Prefix))
	}
	if options.Limit > 0 {
		listOpts = append(listOpts, store.ListLimit(options.Limit))
	}
	if options.Offset > 0 {
		listOpts = append(listOpts, store.ListOffset(options.Offset))
	}

	return DefaultStore.List(listOpts...)
}
