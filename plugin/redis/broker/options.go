package broker

import (
	"context"
	"crypto/tls"

	"github.com/micro/micro/v3/service/broker"
)

type optionsKey struct{}

// Options which are used to configure the redis broker
type Options struct {
	Address   string
	User      string
	Password  string
	TLSConfig *tls.Config
}

// Option is a function which configures options
type Option func(o *Options)

func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

func User(user string) Option {
	return func(o *Options) {
		o.User = user
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func TLSConfig(tlsConfig *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = tlsConfig
	}
}

func RedisOptions(opts Options) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, optionsKey{}, opts)
	}
}
