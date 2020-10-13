module github.com/micro/micro/profile/ci/v3

go 1.15

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201013135405-1a962e46fd3a
	github.com/micro/go-plugins/registry/etcd/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/micro/v3 v3.0.0-beta.4.0.20200921154750-68282c70c194
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/urfave/cli/v2 v2.2.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../..
