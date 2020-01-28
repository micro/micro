package db

type Options struct {
	Url    string
	DBName string
	Table  string
}

type Option func(*Options)

func WithDBName(name string) Option {
	return func(options *Options) {
		options.DBName = name
	}
}

// WithTable set the table to store data, if supported.
func WithTable(table string) Option {
	return func(options *Options) {
		options.Table = table
	}
}

func WithUrl(url string) Option {
	return func(options *Options) {
		options.Url = url
	}
}

type ListOptions struct {
}

type ListOption func(*Options)
