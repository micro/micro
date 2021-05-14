module github.com/micro/micro/plugin/s3/v3

go 1.15

require (
	github.com/aws/aws-sdk-go v1.23.0
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/micro/micro/v3 v3.2.2-0.20210514104957-95ee0dd08833
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0 // indirect
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro/micro/v3 => ../..
