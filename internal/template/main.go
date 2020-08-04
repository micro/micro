package template

var (
	MainSRV = `package main

import (
	"{{.Dir}}/handler"
	{{dehyphen .Alias}} "{{.Dir}}/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Setup service
	srv := service.New(
		service.Name("{{lower .Alias}}"),
		service.Version("latest"),
	)

	// Register Handler
	{{dehyphen .Alias}}.Register{{title .Alias}}Handler(srv.Server(), new(handler.{{title .Alias}}))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
`
)
