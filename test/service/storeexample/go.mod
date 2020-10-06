module example

go 1.13

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201006103121-76291cf85d25
	github.com/micro/micro/v3 v3.0.0-beta.0.20200817234352-e8d00c2dea0d
	google.golang.org/grpc v1.27.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../..
