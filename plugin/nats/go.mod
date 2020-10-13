module github.com/micro/micro/plugin/nats/v3

go 1.13

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201009122815-dad05be95ee0
	github.com/micro/micro/v3 v3.0.0-beta.6
	github.com/nats-io/nats.go v1.10.0
)

replace github.com/micro/micro/v3 => ../..