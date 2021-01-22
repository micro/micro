module github.com/micro/micro/plugin/nats/broker/v3

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/micro/micro/v3 v3.0.4
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
)

replace github.com/micro/micro/v3 => ../../..
