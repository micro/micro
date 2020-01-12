package db

type Options struct {
	Url    string
	DBName string
}

type Option func(*Options)

func WithUrl(url string) Option {
	return func(options *Options) {
		options.Url = url
	}
}

func WithDBName(name string) Option {
	return func(options *Options) {
		options.DBName = name
	}
}
