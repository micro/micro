module example-service

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v3 v3.0.0-alpha.0.20200727163319-c31ce76d407a
	github.com/micro/micro/v3 v2.9.2-0.20200727095830-a9d1f931458a
	google.golang.org/grpc v1.27.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../..
