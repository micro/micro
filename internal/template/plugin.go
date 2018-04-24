package template

var (
	Plugin = `package main
{{if .Plugins}}
import ({{range .Plugins}}
	_ "github.com/micro/go-plugins/{{.}}"{{end}}
){{end}}
`
)
