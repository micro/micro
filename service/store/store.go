package store

import (
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/service"
)

var (
	// DefaultStore implementation
	DefaultStore store.Store = service.NewStore()
)
