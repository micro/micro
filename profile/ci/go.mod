module github.com/micro/micro/profile/ci/v3

go 1.15

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200916120904-3c7f663e8b15
	github.com/micro/go-plugins/registry/etcd/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/micro/v3 v3.0.0-beta.4.0.20200916121553-fe97adbaf454
	github.com/urfave/cli/v2 v2.2.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
