module example

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v3 v3.0.0-alpha.0.20200728125458-9813f98c8b60
	github.com/micro/micro/v3 v3.0.0-20200728090928-ad22505562c9
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../..
