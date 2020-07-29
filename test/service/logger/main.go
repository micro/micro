package main

import (
	"time"

	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
)

func init() {
	go func() {
		for {
			logger.Infof("This is a log line %s\n", time.Now())
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.logger"),
	)
	srv.Run()
}
