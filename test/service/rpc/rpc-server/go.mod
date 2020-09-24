module github.com/micro/micro/v3/test/service/rpc/rpc-server

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/micro/v3 v3.0.0-beta.4.0.20200916121553-fe97adbaf454
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../../..
