package template

var (
	Plugin = `package main
{{if .Plugins}}
import ({{range .Plugins}}
	_ "github.com/micro/go-plugins/v2/{{.}}"{{end}}
){{end}}
`
)
