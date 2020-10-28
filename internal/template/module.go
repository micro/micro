package template

var (
	Module = `module {{.Dir}}

go 1.15

require (
	github.com/micro/micro/v3 v3.0.0-beta.4.0.20200922151713-de8b56c2b15d
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200922112322-927d4f8eced6
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
`
)
