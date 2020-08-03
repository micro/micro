package template

var (
	MainSRV = `package main

import (
	"{{.Dir}}/handler"
	{{dehyphen .Alias}} "{{.Dir}}/proto"

	"github.com/micro/micro/v3/service"
)

func main() {
	// Register Handler
	{{dehyphen .Alias}}.Register{{title .Alias}}Handler(new(handler.{{title .Alias}}))

	// Run service
	service.Run()
}
`
)
