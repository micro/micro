module github.com/micro/micro/cmd/platform

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../..

replace github.com/micro/micro/profile/platform/v3 => ../../profile/platform

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201001155213-a68b7b7b8603 // indirect
	github.com/micro/micro/profile/platform/v3 v3.0.0-20200928084632-c6281c58b123
	github.com/micro/micro/v3 v3.0.0-beta.4.0.20200928084632-c6281c58b123
)
