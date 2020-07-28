module dep-test-service

go 1.13

replace dependency => ../

require (
	dependency v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v3 v3.0.0-20200726113018-34dc21d37191
	google.golang.org/grpc v1.27.0
)

replace github.com/micro/micro/v3 => ../../..