package manager

import "github.com/micro/go-micro/v2/store"

// Options for the runtime manager
type Options struct {
	// Profile contains the env vars to set when
	// running a service
	Profile []string
	// Store to persist state
	Store store.Store
}

// Option sets an option
type Option func(*Options)

// Profile to use when running services
func Profile(p []string) Option {
	return func(o *Options) {
		o.Profile = p
	}
}

// Store to persist services and sync events
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}
