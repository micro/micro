package manager

import "github.com/micro/go-micro/v3/runtime/builder"

type Options struct {
	Builder builder.Builder
}

type Option func(o *Options)

func Builder(b builder.Builder) Option {
	return func(o *Options) {
		o.Builder = b
	}
}
