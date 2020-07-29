package main

import (
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func init() {
	go func() {
		for {
			logger.Infof("This is a log line %s", time.Now())
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
