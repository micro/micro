package store

import (
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/store/client"
)

var (
	// DefaultStore implementation
	DefaultStore store.Store = client.NewStore()
)

type (
	// Record is an alias for store.Record
	Record = store.Record
)

// Read takes a single key name and optional ReadOptions. It returns matching []*Record or an error.
func Read(key string, opts ...store.ReadOption) ([]*Record, error) {
	return DefaultStore.Read(key, opts...)
}

// Write a record to the store, and returns an error if the record was not written.
func Write(r *Record, opts ...store.WriteOption) error {
	return DefaultStore.Write(r, opts...)
}

// Delete removes the record with the corresponding key from the store.
func Delete(key string, opts ...store.DeleteOption) error {
	return DefaultStore.Delete(key, opts...)
}

// List returns any keys that match, or an empty list with no error if none matched.
func List(opts ...store.ListOption) ([]string, error) {
	return DefaultStore.List(opts...)
}
