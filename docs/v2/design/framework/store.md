# Store

Store is an abstraction for key-value storage.

## Overview

For the majority of time microservices are considered stateless and storage is offloaded to a database. 
Considering that we provide a framework, storage and distributed storage needs to be a core concern.
Micro provides a Store interface for key-value storage and a micro store service as the RPC layer 
abstraction.

## Design

The interface is:

```go
// Store is the interface for data storage
type Store interface {
	Init(...Option)                          error
	Options()                                Options
	Read(key string, opts ...ReadOption)     ([]*Record, error)
	Write(*Record, opts ...WriteOption)      error
	Delete(key string, opts ...DeleteOption) error
	List(opts ...ListOption)                 ([]string, error)
	String()                                 string
}

// Record is the data stored by the store
type Record struct {
	// The key for the record
	Key    string
	// The encoded database
	Value  []byte
	// Associated metadata
	Metadata map[string]interface{}
	// Time at which the record expires
	Expiry time.Duration
}
```

### Init

`Init()` initialises the store. It must any required setup on the backing storage
implementation and check that it is ready for use, returning any errors.
`Init()` **must** be called successfully before the store is used.

#### Option

```go
type Options struct {
	Nodes     []string
	Namespace string
	Prefix    string
	Suffix    string
	Context   context.Context
}

type Option func(o *Options)
```

`Nodes` contains the addresses or other connection information of the backing
storage. For example, an etcd implementation would contain the nodes of the
cluster. A SQL implementation could contain one or more connection strings.

`Namespace` allows multiple isolated stores to be kept in one backend, if supported.
For example, multiple tables in a SQL store.

`Prefix` and `Suffix` set a global prefix/suffix on all keys.

`Context` should contain all implementation specific options, using
[`context.WithValue`](https://pkg.go.dev/context?tab=doc#WithValue) as a KV store.

### Options

`Options()` returns the current options

### String

`String()` returns the name of the implementation, e.g. `memory`. Useful for logging purposes.

### Read

`Read()` takes a single key name and optional `ReadOption`s. It returns matching `*Record`s or an error.

#### ReadOption

```go
type ReadOptions struct {
	Prefix bool
	Suffix bool
}

type ReadOption func(r *ReadOptions)
```

`Prefix` and `Suffix` return all keys with matching prefix or suffix.

### Write

`Write()` writes a record to the store, and returns an error if the record was not written.

#### WriteOption

```go
type WriteOptions struct {
  Expiry time.Time
  TTL    time.Duration
}

type WriteOption func(w *WriteOptions)
```

If Expiry or TTL are passed as options, overwrite the record's expiry before writing.

### Delete

`Delete()` removes the record with the corresponding key from the store.

No options are defined yet.

### List

`List()` returns any keys that match, or an empty list with no error if none matched.

#### ListOption

```go
type ListOptions struct {
  Prefix string
  Suffix string
}

type ListOption func(l *ListOptions)
```

If Prefix and / or Suffix are set, the returned list is limited to keys that have the prefix or suffix.

## Caching

Caching is a layer to be built on top of the store in store/cache much like registry/cache.

Caching needs to take into consideration cache coherence and invalidation
- https://en.m.wikipedia.org/wiki/Cache_coherence
- https://en.m.wikipedia.org/wiki/Cache_invalidation

## Indexing

The store supports indexing via metadata field values. These values can be scanned to quickly access 
records that otherwise might of a larger size and more costly to decode. 

In the case of a system like cockroachdb we store metadata as a separate field that uses JSONB format 
so that it can be queried. See reference here https://www.cockroachlabs.com/docs/stable/jsonb.html.

Storage and querying would be of the form

```go

type Fields map[string]string

store.Write(&Record{
	Key: "user:1",
	Value: []byte(...),
	Metadata: map[string]interface{
		"name": "John",
		"email": "john@example.com",
	},
}


wrote.Read("", store.ReadWhere(&store.Fields{
	"name": "john",
})
```
