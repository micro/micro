module github.com/micro/micro/plugin/postgres/v3

go 1.15

require (
	github.com/lib/pq v1.8.0
	github.com/micro/micro/v3 v3.2.2-0.20210526102354-5294ad2ae421
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro/micro/v3 => ../..
