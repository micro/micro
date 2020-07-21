module example-service

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.9.1-0.20200720090451-a3a7434f2cd9
	github.com/micro/micro/v2 v2.9.2-0.20200721134233-06a44ad58f35
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
