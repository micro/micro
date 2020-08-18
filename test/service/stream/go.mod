module github.com/micro/micro/test/service/stream

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../..

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.0.20200817145830-c54d0770dbaa
	github.com/micro/micro/v3 v3.0.0-beta.0.20200817215434-d519cfc25878
	github.com/micro/services v0.10.0 // indirect
)
