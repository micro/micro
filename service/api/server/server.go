package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/api"
	"micro.dev/v4/service"
	"micro.dev/v4/service/api"
	"micro.dev/v4/service/api/auth"
	"micro.dev/v4/service/api/handler"
	aapi "micro.dev/v4/service/api/handler/api"
	"micro.dev/v4/service/api/handler/event"
	ahttp "micro.dev/v4/service/api/handler/http"
	"micro.dev/v4/service/api/handler/rpc"
	"micro.dev/v4/service/api/handler/web"
	"micro.dev/v4/service/api/resolver"
	"micro.dev/v4/service/api/resolver/grpc"
	"micro.dev/v4/service/api/resolver/host"
	"micro.dev/v4/service/api/resolver/path"
	"micro.dev/v4/service/api/resolver/subdomain"
	"micro.dev/v4/service/api/router"
	regRouter "micro.dev/v4/service/api/router/registry"
	httpapi "micro.dev/v4/service/api/server/http"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/registry"
)

var (
	Name      = "api"
	Address   = ":8080"
	Handler   = "meta"
	Resolver  = "micro"
	APIPath   = "/"
	ProxyPath = "/{service:[a-zA-Z0-9]+}"
	Namespace = ""
)

var (
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Set the api address e.g 0.0.0.0:8080",
			EnvVars: []string{"MICRO_API_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "handler",
			Usage:   "Specify the request handler to be used for mapping HTTP requests to services; {api, event, http, rpc}",
			EnvVars: []string{"MICRO_API_HANDLER"},
		},
		&cli.StringFlag{
			Name:    "namespace",
			Usage:   "Set the namespace used by the API e.g. com.example",
			EnvVars: []string{"MICRO_API_NAMESPACE"},
		},
		&cli.StringFlag{
			Name:    "resolver",
			Usage:   "Set the hostname resolver used by the API {host, path, grpc}",
			EnvVars: []string{"MICRO_API_RESOLVER"},
		},
		&cli.BoolFlag{
			Name:    "enable_cors",
			Usage:   "Enable CORS, allowing the API to be called by frontend applications",
			EnvVars: []string{"MICRO_API_ENABLE_CORS"},
			Value:   true,
		},
	}
)

func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("handler")) > 0 {
		Handler = ctx.String("handler")
	}
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}
	if len(ctx.String("api_handler")) > 0 {
		Handler = ctx.String("api_handler")
	}
	if len(ctx.String("api_address")) > 0 {
		Address = ctx.String("api_address")
	}
	// initialise service
	srv := service.New(service.Name(Name))

	// Init API
	var opts []api.Option

	if ctx.Bool("enable_cors") {
		opts = append(opts, api.EnableCORS(true))
	}

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

	// resolver options
	ropts := []resolver.Option{
		resolver.WithServicePrefix(Namespace),
		resolver.WithHandler(Handler),
	}

	// default resolver
	rr := NewResolver(ropts...)

	switch Resolver {
	case "subdomain":
		rr = subdomain.NewResolver(rr)
	case "host":
		rr = host.NewResolver(ropts...)
	case "path":
		rr = path.NewResolver(ropts...)
	case "grpc":
		rr = grpc.NewResolver(ropts...)
	}

	switch Handler {
	case "rpc":
		logger.Infof("Registering API RPC Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(rpc.Handler),
			router.WithResolver(rr),
			router.WithRegistry(registry.DefaultRegistry),
		)
		rp := rpc.NewHandler(
			handler.WithNamespace(Namespace),
			handler.WithRouter(rt),
			handler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(rp)
	case "api":
		logger.Infof("Registering API Request Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(aapi.Handler),
			router.WithResolver(rr),
			router.WithRegistry(registry.DefaultRegistry),
		)
		ap := aapi.NewHandler(
			handler.WithNamespace(Namespace),
			handler.WithRouter(rt),
			handler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(ap)
	case "event":
		logger.Infof("Registering API Event Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(event.Handler),
			router.WithResolver(rr),
			router.WithRegistry(registry.DefaultRegistry),
		)
		ev := event.NewHandler(
			handler.WithNamespace(Namespace),
			handler.WithRouter(rt),
			handler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(ev)
	case "http":
		logger.Infof("Registering API HTTP Handler at %s", ProxyPath)
		rt := regRouter.NewRouter(
			router.WithHandler(ahttp.Handler),
			router.WithResolver(rr),
			router.WithRegistry(registry.DefaultRegistry),
		)
		ht := ahttp.NewHandler(
			handler.WithNamespace(Namespace),
			handler.WithRouter(rt),
			handler.WithClient(srv.Client()),
		)
		r.PathPrefix(ProxyPath).Handler(ht)
	case "web":
		logger.Infof("Registering API Web Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(web.Handler),
			router.WithResolver(rr),
			router.WithRegistry(registry.DefaultRegistry),
		)
		w := web.NewHandler(
			handler.WithNamespace(Namespace),
			handler.WithRouter(rt),
			handler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(w)
	default:
		logger.Infof("Registering API Default Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithResolver(rr),
			router.WithRegistry(registry.DefaultRegistry),
		)
		r.PathPrefix(APIPath).Handler(Meta(srv.Client(), rt, Namespace))
	}

	// append the auth wrapper
	h = auth.Wrapper(rr, Namespace)(h)

	// create a new api server with wrappers
	api := httpapi.NewServer(Address)
	// initialise
	api.Init(opts...)
	// register the http handler
	api.Handle("/", h)

	// Start API
	if err := api.Start(); err != nil {
		logger.Fatal(err)
	}

	// register the rpc handler
	pb.RegisterApiHandler(srv.Server(), &handler.APIHandler{})

	// Run server
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		logger.Fatal(err)
	}

	return nil
}
