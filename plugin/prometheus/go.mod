module github.com/micro/micro/plugin/prometheus/v3

go 1.15

require (
	github.com/micro/micro/v3 v3.0.4
	github.com/prometheus/client_golang v1.11.1
	github.com/stretchr/testify v1.8.0
)

replace github.com/micro/micro/v3 => ../..
