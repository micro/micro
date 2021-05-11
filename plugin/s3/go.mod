module github.com/micro/micro/plugin/s3/v3

go 1.15

require (
	github.com/micro/micro/v3 v3.2.1
	github.com/minio/minio-go/v7 v7.0.10
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro/micro/v3 => ../..

