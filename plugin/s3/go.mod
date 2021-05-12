module github.com/micro/micro/plugin/s3/v3

go 1.15

require (
	github.com/aws/aws-sdk-go v1.23.0
	github.com/micro/micro/v3 v3.2.2-0.20210512111111-9e3e99a26e10
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro/micro/v3 => ../..
