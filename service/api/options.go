package api

type Options struct{}

type Option func(*Options) error
