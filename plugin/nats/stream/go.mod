module github.com/micro/micro/plugin/nats/stream/v3

go 1.15

require (
	github.com/google/uuid v1.1.2
	github.com/micro/micro/v3 v3.0.4
	github.com/nats-io/nats-streaming-server v0.24.3 // indirect
	github.com/nats-io/nats.go v1.13.1-0.20220308171302-2f2f6968e98d
	github.com/nats-io/stan.go v0.10.2
	github.com/pkg/errors v0.9.1
)

replace github.com/micro/micro/v3 => ../../..
