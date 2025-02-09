package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/micro/micro/v5/service/api"
	"github.com/micro/micro/v5/service/api/handler"
	"github.com/micro/micro/v5/service/api/handler/rpc"
	"github.com/micro/micro/v5/service/api/router"
	regRouter "github.com/micro/micro/v5/service/api/router/registry"
	httpapi "github.com/micro/micro/v5/service/api/server/http"
	"github.com/micro/micro/v5/service/client"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/registry"
	"github.com/urfave/cli/v2"
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
