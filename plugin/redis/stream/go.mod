module github.com/micro/micro/plugin/redis/stream/v3

go 1.15

require (
	github.com/go-redis/redis/v8 v8.10.1-0.20210615084835-43ec1464d9a6
	github.com/google/uuid v1.1.2
	github.com/micro/micro/v3 v3.2.2-0.20210607125937-d152c4ec7749
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro/micro/v3 => ../../..
