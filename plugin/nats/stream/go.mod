module github.com/micro/micro/plugin/nats/stream/v3

go 1.15

require (
	github.com/google/uuid v1.1.2
	github.com/micro/micro/v3 v3.0.4
	github.com/nats-io/nats-streaming-server v0.19.0 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/stan.go v0.7.0
	github.com/pkg/errors v0.9.1
)

replace github.com/micro/micro/v3 => ../../..
