package template

var (
	Module = `module {{.Dir}}

go {{.GoVersion}}

require github.com/micro/micro/v5 latest
`
)
