module logger

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../..

require (
	github.com/micro/go-micro/v3 v3.0.0-alpha.0.20200728125458-9813f98c8b60
	github.com/micro/micro/v3 v3.0.0-00010101000000-000000000000
)
