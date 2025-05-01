package main

import (
	"log"

	httpServer "github.com/micro/plugins/v5/server/http"
	"go-micro.dev/v5"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/server"
)

const (
	SERVER_NAME = "demo-http" // server name
)

func main() {

	srv := httpServer.NewServer(
		server.Name(SERVER_NAME),
		server.Address(":8080"),
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// register router
	demo := newDemo()
	demo.InitRouter(router)

	hd := srv.NewHandler(router)
	if err := srv.Handle(hd); err != nil {
		log.Fatalln(err)
	}

	service := micro.NewService(
		micro.Server(srv),
		micro.Registry(registry.NewRegistry()),
	)
	service.Init()
	service.Run()
}

// demo
type demo struct{}

func newDemo() *demo {
	return &demo{}
}

func (a *demo) InitRouter(router *gin.Engine) {
	router.POST("/demo", a.demo)
}

func (a *demo) demo(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "call go-micro v3 http server success"})
}
