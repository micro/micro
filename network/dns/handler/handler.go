// Package handler is the handler for the `micro network dns` command
package handler

import (
	"github.com/micro/micro/network/dns/provider/cloudflare"
)

// New returns a new handler
func New() *DNS {
	provider, _ := cloudflare.New()
	return &DNS{
		provider: provider,
	}
}
