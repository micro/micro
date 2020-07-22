package store

import (
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/micro/v2/service/store/client"
)

var (
	// DefaultStore implementation
	DefaultStore store.Store = client.NewStore()
)
