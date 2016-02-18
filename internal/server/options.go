package server

import (
	"crypto/tls"
)

type Option func(o *Options)

type Options struct {
	EnableTLS bool
	TLSConfig *tls.Config
}

func EnableTLS(b bool) Option {
	return func(o *Options) {
		o.EnableTLS = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}
