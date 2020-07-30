package template

var (
	MainSRV = `package main

import (
	"github.com/micro/micro/v3/service"
	
	"{{.Dir}}/handler"
	"{{.Dir}}/subscriber"
	{{dehyphen .Alias}} "{{.Dir}}/proto"
)

func main() {
	// Register Handler
	{{dehyphen .Alias}}.Register{{title .Alias}}Handler(new(handler.{{title .Alias}}))

	// Register Struct as Subscriber
	service.RegisterSubscriber("{{.Alias}}", new(subscriber.{{title .Alias}}))

	// Run service
	service.Run()
}
`
)
