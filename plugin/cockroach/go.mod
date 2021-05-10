module github.com/micro/micro/plugin/cockroach/v3

go 1.15

require (
	github.com/lib/pq v1.8.0
	github.com/micro/micro/v3 v3.0.4
	github.com/pkg/errors v0.9.1
)

replace github.com/micro/micro/v3 => ../..
