package template

var (
	MainSRV = `package main

import (
	"github.com/micro/go-log"

	"github.com/micro/go-micro"
	"{{.Dir}}/handler"
	"{{.Dir}}/subscriber"

	example "{{.Dir}}/proto/example"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("{{.FQDN}}"),
		micro.Version("latest"),
	)

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Register Struct as Subscriber
	service.Server().Subscribe(
		service.Server().NewSubscriber("topic.{{.FQDN}}", new(subscriber.Example)),
	)

	// Register Function as Subscriber
	service.Server().Subscribe(
		service.Server().NewSubscriber("topic.{{.FQDN}}", subscriber.Handler),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
`
	MainAPI = `package main

import (
	"github.com/micro/go-log"

	"github.com/micro/go-micro"
	"{{.Dir}}/handler"
	"{{.Dir}}/client"

	example "{{.Dir}}/proto/example"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("{{.FQDN}}"),
		micro.Version("latest"),
	)

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Initialise service
	service.Init(
		// create wrap for the Example srv client
		micro.WrapHandler(client.ExampleWrapper(service)),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
`
	MainWEB = `package main

import (
        "github.com/micro/go-log"
	"net/http"

        "github.com/micro/go-web"
        "{{.Dir}}/handler"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("{{.FQDN}}"),
                web.Version("latest"),
        )

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	service.HandleFunc("/example/call", handler.ExampleCall)

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
`
)
