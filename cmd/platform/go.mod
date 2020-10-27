module github.com/micro/micro/cmd/platform

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../..

replace github.com/micro/micro/profile/platform/v3 => ../../profile/platform

replace github.com/micro/micro/plugin/etcd/v3 => ../../plugin/etcd

replace github.com/micro/micro/plugin/cockroach/v3 => ../../plugin/cockroach

replace github.com/micro/micro/plugin/prometheus/v3 => ../../plugin/prometheus

replace github.com/micro/micro/plugin/nats/broker/v3 => ../../plugin/nats/broker

replace github.com/micro/micro/plugin/nats/stream/v3 => ../../plugin/nats/stream

require (
	github.com/micro/micro/profile/platform/v3 v3.0.0-20200928084632-c6281c58b123
	github.com/micro/micro/v3 v3.0.0-beta.6
)
