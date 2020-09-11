module github.com/micro/micro/profile/ci/v3

go 1.15

require (
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200908121445-a02a25d9558e
	github.com/micro/go-plugins/registry/etcd/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/micro/v3 v3.0.0-beta.3.0.20200908134309-90be716874c4
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
