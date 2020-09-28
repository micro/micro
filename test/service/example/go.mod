module example

go 1.13

require (
	github.com/bwmarrin/discordgo v0.20.2 // indirect
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/go-git/go-git/v5 v5.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/lucas-clemente/quic-go v0.14.1 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200917131714-7750f542b4c2
	github.com/micro/micro/v3 v3.0.0-20200728090928-ad22505562c9
	github.com/nats-io/nats-streaming-server v0.18.0 // indirect
	github.com/nlopes/slack v0.6.1-0.20191106133607-d06c2a2b3249 // indirect
	github.com/prometheus/client_golang v1.7.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/tools v0.0.0-20200117065230-39095c1d176c // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../..
