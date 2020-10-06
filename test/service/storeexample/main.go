package main

import (
	"example/handler"

	"github.com/micro/micro/v3/service"
)

func main() {
	srv := service.New(service.Name("storeexample"))

	srv.Handle(new(handler.Example))

	srv.Run()
}
