package main

import (
	"time"

	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v2/service"
)

func init() {
	go func() {
		for {
			logger.Info("These logs will happen until you stop me! Never stop never stopping!")
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.logspammer"),
		service.Version("latest"),
	)
	srv.Run()
}
