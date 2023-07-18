package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"micro.dev/v4/service/api"
	"micro.dev/v4/service/api/handler"
	"micro.dev/v4/service/api/handler/rpc"
	"micro.dev/v4/service/api/router"
	regRouter "micro.dev/v4/service/api/router/registry"
	httpapi "micro.dev/v4/service/api/server/http"
	"micro.dev/v4/service/client"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/registry"
)

func runAPI(ctx *cli.Context, wait chan bool) error {
	// Init API
	var opts []api.Option

	opts = append(opts, api.EnableCORS(true))

	// create the router
	var h http.Handler
	r := mux.NewRouter()
	h = r

	// return version and list of services
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		response := fmt.Sprintf(`{"version": "%s"}`, ctx.App.Version)
		w.Write([]byte(response))
	})

	// strip favicon.ico
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	rt := regRouter.NewRouter(router.WithRegistry(registry.DefaultRegistry))
	r.PathPrefix("/").Handler(rpc.NewHandler(
		handler.WithClient(client.DefaultClient),
		handler.WithRouter(rt),
	))

	// create a new api server with wrappers
	api := httpapi.NewServer(":8080")
	// initialise
	api.Init(opts...)
	// register the http handler
	api.Handle("/", authWrapper()(h))

	// Start API
	if err := api.Start(); err != nil {
		logger.Fatal(err)
	}

	// wait to stop
	<-wait

	// Stop API
	if err := api.Stop(); err != nil {
		logger.Fatal(err)
	}

	return nil
}
