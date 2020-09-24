module github.com/micro/micro/profile/platform/v3

go 1.15

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200921154545-9dbd75f2cc13
	github.com/micro/go-plugins/broker/nats/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/go-plugins/events/stream/nats/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/go-plugins/metrics/prometheus/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/go-plugins/registry/etcd/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/go-plugins/store/cockroach/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/micro/profile/ci/v3 v3.0.0-20200918083418-64086487dabb // indirect
	github.com/micro/micro/v3 v3.0.0-beta.4.0.20200921154750-68282c70c194
	github.com/urfave/cli/v2 v2.2.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
