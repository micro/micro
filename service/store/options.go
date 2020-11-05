// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/store/options.go

package store

import "context"

// Options to use when reading from the store
type Options struct {
	// Prefix scopes the query to records that are prefixed with key
	Prefix string
	// Limit limits the number of returned records
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// Option sets values in Options
type Option func(o *Options)

// Prefix returns all records that are prefixed with key
func Prefix(p string) Option {
	return func(r *Options) {
		r.Prefix = p
	}
}

// Limit limits the number of responses to l
func Limit(l uint) Option {
	return func(r *Options) {
		r.Limit = l
	}
}

// Offset starts returning responses from o. Use in conjunction with Limit for pagination
func Offset(o uint) Option {
	return func(r *Options) {
		r.Offset = o
	}
}

// StoreOptions contains configuration for the Store
type StoreOptions struct {
	// Nodes contains the addresses or other connection information of the backing storage.
	// For example, an etcd implementation would contain the nodes of the cluster.
	// A SQL implementation could contain one or more connection strings.
	Nodes []string
	// Database allows multiple isolated stores to be kept in one backend, if supported.
	Database string
	// Table is analagous to a table in database backends or a key prefix in KV backends
	Table string
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context
}

// StoreOption sets values in StoreOptions
type StoreOption func(o *StoreOptions)

// Nodes contains the addresses or other connection information of the backing storage.
// For example, an etcd implementation would contain the nodes of the cluster.
// A SQL implementation could contain one or more connection strings.
func Nodes(a ...string) StoreOption {
	return func(o *StoreOptions) {
		o.Nodes = a
	}
}

// Database allows multiple isolated stores to be kept in one backend, if supported.
func Database(db string) StoreOption {
	return func(o *StoreOptions) {
		o.Database = db
	}
}

// Table is analagous to a table in database backends or a key prefix in KV backends
func Table(t string) StoreOption {
	return func(o *StoreOptions) {
		o.Table = t
	}
}

// WithContext sets the stores context, for any extra configuration
func WithContext(c context.Context) StoreOption {
	return func(o *StoreOptions) {
		o.Context = c
	}
}

// ReadOptions configures an individual Read operation
type ReadOptions struct {
	Database, Table string
	// Prefix returns all records that are prefixed with key
	Prefix bool
	// Suffix returns all records that have the suffix key
	Suffix bool
	// Limit limits the number of returned records
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// ReadOption sets values in ReadOptions
type ReadOption func(r *ReadOptions)

// ReadFrom the database and table
func ReadFrom(database, table string) ReadOption {
	return func(r *ReadOptions) {
		r.Database = database
		r.Table = table
	}
}

// ReadPrefix returns all records that are prefixed with key
func ReadPrefix() ReadOption {
	return func(r *ReadOptions) {
		r.Prefix = true
	}
}

// ReadSuffix returns all records that have the suffix key
func ReadSuffix() ReadOption {
	return func(r *ReadOptions) {
		r.Suffix = true
	}
}

// ReadLimit limits the number of responses to l
func ReadLimit(l uint) ReadOption {
	return func(r *ReadOptions) {
		r.Limit = l
	}
}

// ReadOffset starts returning responses from o. Use in conjunction with Limit for pagination
func ReadOffset(o uint) ReadOption {
	return func(r *ReadOptions) {
		r.Offset = o
	}
}

// WriteOptions configures an individual Write operation
// If Expiry and TTL are set TTL takes precedence
type WriteOptions struct {
	Database, Table string
}

// WriteOption sets values in WriteOptions
type WriteOption func(w *WriteOptions)

// WriteTo the database and table
func WriteTo(database, table string) WriteOption {
	return func(w *WriteOptions) {
		w.Database = database
		w.Table = table
	}
}

// DeleteOptions configures an individual Delete operation
type DeleteOptions struct {
	Database, Table string
}

// DeleteOption sets values in DeleteOptions
type DeleteOption func(d *DeleteOptions)

// DeleteFrom the database and table
func DeleteFrom(database, table string) DeleteOption {
	return func(d *DeleteOptions) {
		d.Database = database
		d.Table = table
	}
}

// ListOptions configures an individual List operation
type ListOptions struct {
	// List from the following
	Database, Table string
	// Prefix returns all keys that are prefixed with key
	Prefix string
	// Suffix returns all keys that end with key
	Suffix string
	// Limit limits the number of returned keys
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// ListOption sets values in ListOptions
type ListOption func(l *ListOptions)

// ListFrom the database and table
func ListFrom(database, table string) ListOption {
	return func(l *ListOptions) {
		l.Database = database
		l.Table = table
	}
}

// ListPrefix returns all keys that are prefixed with key
func ListPrefix(p string) ListOption {
	return func(l *ListOptions) {
		l.Prefix = p
	}
}

// ListSuffix returns all keys that end with key
func ListSuffix(s string) ListOption {
	return func(l *ListOptions) {
		l.Suffix = s
	}
}

// ListLimit limits the number of returned keys to l
func ListLimit(l uint) ListOption {
	return func(lo *ListOptions) {
		lo.Limit = l
	}
}

// ListOffset starts returning responses from o. Use in conjunction with Limit for pagination.
func ListOffset(o uint) ListOption {
	return func(l *ListOptions) {
		l.Offset = o
	}
}
