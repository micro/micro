module github.com/micro/micro/plugin/s3/v3

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.0
	github.com/micro/micro/v3 v3.3.1-0.20210803122146-2a2fa437600d
	github.com/stretchr/testify v1.8.0
)

replace github.com/micro/micro/v3 => ../..
