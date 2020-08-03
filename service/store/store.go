package store

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/internal/cmd"
)

func init() {
	cmd.Init(func(ctx *cli.Context) error {
		var opts []store.Option
		if len(ctx.String("store_address")) > 0 {
			opts = append(opts, store.Nodes(strings.Split(ctx.String("store_address"), ",")...))
		}
		if len(ctx.String("namespace")) > 0 {
			opts = append(opts, store.Database(ctx.String("namespace")))
		}
		return DefaultStore.Init(opts...)
	})
}

var (
	// DefaultStore implementation
	DefaultStore store.Store
)

// Read takes a single key name and optional ReadOptions. It returns matching []*Record or an error.
func Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	return DefaultStore.Read(key, opts...)
}

// Write a record to the store, and returns an error if the record was not written.
func Write(r *store.Record, opts ...store.WriteOption) error {
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
