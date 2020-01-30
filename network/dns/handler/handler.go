// Package handler is the handler for the `micro network dns` command
package handler

import (
	"github.com/micro/micro/v2/network/dns/provider"
)

// New returns a new handler
func New(provider provider.Provider, token string) *DNS {
	return &DNS{
		provider:    provider,
		bearerToken: token,
	}
}
