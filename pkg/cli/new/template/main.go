package template

var (
	MainSRV = `package main

import (
	"{{.Dir}}/handler"
	pb "{{.Dir}}/proto"

	"go-micro.dev/v5"
)

func main() {
	// Create service
	service := micro.New("{{lower .Alias}}")

	// Register handler
	pb.Register{{title .Alias}}Handler(service.Server(), handler.New())

	// Run service
	service.Run()
}
`
)
