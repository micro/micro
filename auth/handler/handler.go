package handler

import (
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/basic"
	"github.com/micro/go-micro/v2/store"
	memStore "github.com/micro/go-micro/v2/store/memory"
)

var joinKey = ":"

// Handler processes RPC calls
type Handler struct {
	Options        auth.Options
	SecretProvider token.Provider
	TokenProvider  token.Provider
}

// Init the auth
func (h *Handler) Init(opts ...auth.Option) {
	for _, o := range opts {
		o(&h.Options)
	}

	// use the default store as a fallback
	if h.Options.Store == nil {
		h.Options.Store = store.DefaultStore
	}

	// noop will not work for auth
	if h.Options.Store.String() == "noop" {
		h.Options.Store = memStore.NewStore()
	}

	if h.TokenProvider == nil {
		h.TokenProvider = basic.NewTokenProvider(token.WithStore(h.Options.Store))
	}
	if h.SecretProvider == nil {
		h.SecretProvider = basic.NewTokenProvider(token.WithStore(h.Options.Store))
	}
}
