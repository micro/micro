package store

import (
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/store/client"
)

var (
	// DefaultStore implementation
	DefaultStore store.Store = client.NewStore()
)
