module github.com/micro/micro/profile/platform/v3

go 1.15

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201009122815-dad05be95ee0
	github.com/micro/go-plugins/events/stream/nats/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/go-plugins/metrics/prometheus/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/go-plugins/registry/etcd/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/go-plugins/store/cockroach/v3 v3.0.0-20200908121001-4ea6f6760baf
	github.com/micro/micro/plugin/nats/v3 v3.0.0-00010101000000-000000000000
	github.com/micro/micro/v3 v3.0.0-beta.6
	github.com/urfave/cli/v2 v2.2.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../..

replace github.com/micro/micro/plugin/nats/v3 => ../../plugin/nats
