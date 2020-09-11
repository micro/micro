package store

// Options to use when reading from the store
type Options struct {
	// Prefix returns all records that are prefixed with key
	Prefix bool
	// Limit limits the number of returned records
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// Option sets values in Options
type Option func(o *Options)

// Prefix returns all records that are prefixed with key
func Prefix() Option {
	return func(r *Options) {
		r.Prefix = true
	}
}

// Limit limits the number of responses to l
func Limit(l uint) Option {
	return func(r *Options) {
		r.Limit = l
	}
}

// Offset starts returning responses from o. Use in conjunction with Limit for pagination
func Offset(o uint) Option {
	return func(r *Options) {
		r.Offset = o
	}
}
