module github.com/micro/micro/profile/ci/v3

go 1.15

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200918112555-9168c7c61064
	github.com/micro/go-plugins/registry/etcd/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/micro/v3 v3.0.0-beta.3.0.20200913073948-c44509882f8c
	github.com/urfave/cli/v2 v2.2.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
