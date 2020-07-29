package main

import (
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	go func() {
		for {
			logger.Infof("This is a log line %s", time.Now())
			time.Sleep(5 * time.Second)
		}
	}()

	service.Run()
}
