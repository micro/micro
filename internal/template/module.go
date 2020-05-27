package template

var (
	Module = `module {{.Dir}}

	go 1.13

	require (
		github.com/golang/protobuf v1.4.2
		github.com/micro/go-micro/v2 v2.7.0
		go.etcd.io/etcd v0.5.0-alpha.5.0.20200306183522-221f0cc107cb
		google.golang.org/protobuf v1.24.0
	)

	// @todo remove this replace, info: https://github.com/etcd-io/etcd/pull/11564
	replace google.golang.org/grpc => google.golang.org/grpc v1.26.0`
)
