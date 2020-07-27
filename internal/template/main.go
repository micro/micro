package template

var (
	MainFNC = `package main

import (
  log	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"{{.Dir}}/handler"
	"{{.Dir}}/subscriber"
)

func main() {
	// New Service
	function := micro.NewFunction(
		micro.Name("{{.FQDN}}"),
	)

	// Initialise function
	function.Init()

	// Register Handler
	function.Handle(new(handler.{{title .Alias}}))

	// Register Struct as Subscriber
	function.Subscribe("{{.FQDN}}", new(subscriber.{{title .Alias}}))

	// Run service
	if err := function.Run(); err != nil {
		log.Fatal(err)
	}
}
`

	MainSRV = `package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"{{.Dir}}/handler"
	"{{.Dir}}/subscriber"

	{{dehyphen .Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("{{.FQDN}}"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	{{dehyphen .Alias}}.Register{{title .Alias}}Handler(service.Server(), new(handler.{{title .Alias}}))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("{{.FQDN}}", service.Server(), new(subscriber.{{title .Alias}}))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
`
	MainAPI = `package main

import (
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/go-micro/v2"
	"{{.Dir}}/handler"
	"{{.Dir}}/client"

	{{dehyphen .Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("{{.FQDN}}"),
	)

	// Initialise service
	service.Init(
		// create wrap for the {{title .Alias}} service client
		micro.WrapHandler(client.{{title .Alias}}Wrapper(service)),
	)

	// Register Handler
	{{dehyphen .Alias}}.Register{{title .Alias}}Handler(service.Server(), new(handler.{{title .Alias}}))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
`
	MainWEB = `package main

import (
        log "github.com/micro/go-micro/v2/logger"
	      "net/http"
        "github.com/micro/go-micro/v2/web"
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
