package template

var (
	Module = `module {{.Dir}}

go {{.GoVersion}}

require (
	micro.dev/v4 latest
	google.golang.org/protobuf latest
	google.golang.org/grpc latest
)
`
)
