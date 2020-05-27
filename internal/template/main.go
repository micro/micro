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
		micro.Version("latest"),
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

	GoModFNC = `module {{ title .Alias }}

	go 1.13

	require (
		github.com/golang/protobuf v1.4.2
		github.com/micro/go-micro/v2 v2.7.0
		go.etcd.io/etcd v0.5.0-alpha.5.0.20200306183522-221f0cc107cb
		google.golang.org/protobuf v1.24.0
	)

	// @todo remove this replace, info: https://github.com/etcd-io/etcd/pull/11564
	replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
`

	MainSRV = `package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"{{.Dir}}/handler"
	"{{.Dir}}/subscriber"

	{{.Alias}} "{{.Dir}}/proto/{{.Alias}}"
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
	{{.Alias}}.Register{{title .Alias}}Handler(service.Server(), new(handler.{{title .Alias}}))

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

	{{.Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("{{.FQDN}}"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the {{title .Alias}} service client
		micro.WrapHandler(client.{{title .Alias}}Wrapper(service)),
	)

	// Register Handler
	{{.Alias}}.Register{{title .Alias}}Handler(service.Server(), new(handler.{{title .Alias}}))

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
                web.Version("latest"),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	service.HandleFunc("/{{.Alias}}/call", handler.{{title .Alias}}Call)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
`
)
