package template

var (
	Module = `module {{.Dir}}

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
`
)
