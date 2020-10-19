module github.com/micro/micro/test/vendortest/v3

go 1.13

replace github.com/micro/micro/v3 => ../..

replace github.com/micro/test => ../../../test

require (
	github.com/micro/micro/v3 v3.0.0-beta.6
	github.com/micro/test v0.0.0-20201019141839-d620c1e4fc1c
)

replace github.com/micro/micro/plugin/etcd/v3 => ../../plugin/etcd

replace github.com/micro/micro/profile/ci/v3 => ../../profile/ci

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0