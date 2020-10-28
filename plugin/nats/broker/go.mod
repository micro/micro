module github.com/micro/micro/plugin/nats/broker/v3

go 1.15

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201013135405-1a962e46fd3a
	github.com/micro/micro/v3 v3.0.0-beta.6
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
)

replace github.com/micro/micro/v3 => ../../..
