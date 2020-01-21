package db

type Options struct {
	Url    string
	DBName string
}

type Option func(*Options)

func WithDBName(name string) Option {
	return func(options *Options) {
		options.DBName = name
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
