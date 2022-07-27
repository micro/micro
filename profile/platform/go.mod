module github.com/micro/micro/profile/platform

go 1.15

require (
	github.com/go-redis/redis/v8 v8.11.3
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/micro/micro/plugin/etcd/v3 v3.0.0-20201217215412-2f7ad18595ff
	github.com/micro/micro/plugin/postgres/v3 v3.0.0-20210825142032-d27318700a59
	github.com/micro/micro/plugin/prometheus/v3 v3.0.0-20201217215412-2f7ad18595ff
	github.com/micro/micro/plugin/redis/blocklist/v3 v3.0.0-20210818101332-9a333769a3d7
	github.com/micro/micro/plugin/redis/broker/v3 v3.0.0-20210622092058-b12502169757
	github.com/micro/micro/plugin/redis/stream/v3 v3.0.0-20210622092058-b12502169757
	github.com/micro/micro/plugin/s3/v3 v3.0.0-20220203170347-2c86e5fb6d13
	github.com/micro/micro/v3 v3.11.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
