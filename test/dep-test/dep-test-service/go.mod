module dep-test-service

go 1.15

replace dependency => ../

require (
	dependency v0.0.0-00010101000000-000000000000
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/bwmarrin/discordgo v0.20.2 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/forestgiant/sliceutil v0.0.0-20160425183142-94783f95db6c // indirect
	github.com/go-git/go-git/v5 v5.1.0 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/lucas-clemente/quic-go v0.14.1 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201009122815-dad05be95ee0
	github.com/micro/micro/v3 v3.0.0-20200728090928-ad22505562c9
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/nats-io/nats-streaming-server v0.18.0 // indirect
	github.com/nlopes/slack v0.6.1-0.20191106133607-d06c2a2b3249 // indirect
	github.com/prometheus/client_golang v1.7.0 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/tools v0.0.0-20200117065230-39095c1d176c // indirect
	google.golang.org/grpc v1.27.0
	gopkg.in/telegram-bot-api.v4 v4.6.4 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace github.com/micro/micro/v3 => ../../..

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
