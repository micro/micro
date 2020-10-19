module github.com/micro/micro/profile/ci/v3

go 1.15

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/micro/go-micro v1.18.0 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201013135405-1a962e46fd3a
	github.com/micro/micro/plugin/etcd/v3 v3.0.0-00010101000000-000000000000
	github.com/micro/micro/v3 v3.0.0-beta.6
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/urfave/cli/v2 v2.2.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/plugin/etcd/v3 => ../../plugin/etcd

replace github.com/micro/micro/v3 => ../..
