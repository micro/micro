module github.com/micro/micro/cmd/platform

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
replace github.com/micro/micro/v3 => ../..
replace github.com/micro/micro/profile/platform/v3 => ../../profile/platform

require (
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200926122909-017e156440aa
	github.com/micro/micro/profile/platform/v3 v3.0.0-20200928084632-c6281c58b123
	github.com/micro/micro/v3 v3.0.0-beta.4.0.20200928084632-c6281c58b123
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/nats-io/nats-streaming-server v0.18.0 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	golang.org/x/tools v0.0.0-20200117065230-39095c1d176c // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)
