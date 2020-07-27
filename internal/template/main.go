package template

var (
	MainSRV = `package main

import (
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v2/service"
	"{{.Dir}}/handler"
	"{{.Dir}}/subscriber"

	{{dehyphen .Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("{{.FQDN}}"),
		service.Version("latest"),
	)

	// Initialise service
	srv.Init()

	// Register Handler
	{{dehyphen .Alias}}.Register{{title .Alias}}Handler(srv.Server(), new(handler.{{title .Alias}}))

	// Register Struct as Subscriber
	service.RegisterSubscriber("{{.FQDN}}", new(subscriber.{{title .Alias}}))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
`
	MainAPI = `package main

import (
	log "github.com/micro/go-micro/v3/logger"

	"github.com/micro/micro/v2/service"
	"{{.Dir}}/handler"
	"{{.Dir}}/client"

	{{dehyphen .Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("{{.FQDN}}"),
	)

	// Initialise service
	srv.Init(
		// create wrap for the {{title .Alias}} srv client
		micro.WrapHandler(client.{{title .Alias}}Wrapper(srv)),
	)

	// Register Handler
	{{dehyphen .Alias}}.Register{{title .Alias}}Handler(srv.Server(), new(handler.{{title .Alias}}))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
`
	MainWEB = `package main

import (
        log "github.com/micro/go-micro/v3/logger"
	      "net/http"
        "github.com/micro/go-micro/v3/web"
        "{{.Dir}}/handler"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("{{.FQDN}}"),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	service.HandleFunc("/{{dehyphen .Alias}}/call", handler.{{title .Alias}}Call)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
`
)
