module github.com/micro/micro/plugin/redis/stream/v3

go 1.15

require (
	github.com/go-redis/redis/v8 v8.8.3
	github.com/m3o/platform/profile/platform v0.0.0-20210514113832-baec7d77b8f6 // indirect
	github.com/micro/micro/v3 v3.2.2-0.20210515174306-b0144d41f782
)

replace github.com/micro/micro/v3 => ../../..
