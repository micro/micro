package store

// Options to use when reading from the store
type Options struct {
	// Prefix scopes the query to records that are prefixed with key
	Prefix string
	// Limit limits the number of returned records
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// Option sets values in Options
type Option func(o *Options)

// Prefix returns all records that are prefixed with key
func Prefix(p string) Option {
	return func(r *Options) {
		r.Prefix = p
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
