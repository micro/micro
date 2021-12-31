module dep-test-service

go 1.15

replace dependency => ../

require (
	dependency v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.3
	github.com/micro/micro/v3 v3.0.0-20200728090928-ad22505562c9
	google.golang.org/grpc v1.27.0
)

replace github.com/micro/micro/v3 => ../../..

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
