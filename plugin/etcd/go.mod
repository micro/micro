module github.com/micro/micro/plugin/etcd/v3

go 1.15

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201013135405-1a962e46fd3a
	github.com/micro/micro/v3 v3.0.0-beta.6
	github.com/mitchellh/hashstructure v1.0.0
	go.uber.org/zap v1.16.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../..
