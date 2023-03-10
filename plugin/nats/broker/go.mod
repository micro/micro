module github.com/micro/micro/plugin/nats/broker/v3

go 1.15

require (
	github.com/golang/protobuf v1.5.2
	github.com/micro/micro/v3 v3.0.4
	github.com/nats-io/nats-server/v2 v2.7.4 // indirect
	github.com/nats-io/nats.go v1.13.1-0.20220308171302-2f2f6968e98d
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
)

replace github.com/micro/micro/v3 => ../../..
