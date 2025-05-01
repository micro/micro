package main

import (
	"time"

	"go-micro.dev/v5"
	"go-micro.dev/v5/logger"
)

func main() {
	service := micro.NewService(
		micro.Name("com.example.srv.foo"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)
	service.Init()

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
