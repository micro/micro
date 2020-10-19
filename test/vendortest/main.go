package main

import (
	"fmt"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/test/vendortest"
)

func main() {
	fmt.Println("foo:", vendortest.Foo())

	srv := service.New(
		service.Name("vendortest"),
	)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
