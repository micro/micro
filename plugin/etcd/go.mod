module github.com/micro/micro/plugin/etcd/v3

go 1.20

require (
	github.com/micro/micro/v3 v3.0.4
	github.com/mitchellh/hashstructure v1.0.0
	go.etcd.io/etcd/api/v3 v3.5.9
	go.etcd.io/etcd/client/v3 v3.5.9
	go.uber.org/zap v1.17.0
)

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.9 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.54.1 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/soheilhy/cmux => github.com/soheilhy/cmux v0.1.5

replace github.com/micro/micro/v3 => ../..

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
