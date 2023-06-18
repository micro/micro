module github.com/micro/micro/plugin/redis/broker/v3

go 1.15

require (
	github.com/go-redis/redis/v8 v8.10.1-0.20210615084835-43ec1464d9a6
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.0
	github.com/micro/micro/v3 v3.18.6-0.20230605211001-e67a4087d2d2
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
)

replace github.com/micro/micro/v3 => ../../..
