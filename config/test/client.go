package main

import (
	"github.com/micro/go-micro/client/grpc"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/service"
	"github.com/micro/go-micro/util/log"
)

func main() {
	service.DefaultClient = grpc.NewClient()
	src := service.NewSource(
		service.ServiceName("go.micro.config"),
		service.Key("NAMESPACE:CONFIG"),
		service.Path(""),
	)
	if err := config.Load(src); err != nil {
		panic(err)
	}

	v := config.Get("a")
	log.Info(string(v.Bytes()))
}
