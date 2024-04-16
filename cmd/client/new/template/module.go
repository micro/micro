package template

var (
	Module = `module {{.Dir}}

go {{.GoVersion}}

require micro.dev/v4 latest
`
)
