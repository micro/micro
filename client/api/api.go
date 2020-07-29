// Package api is an API Gateway
package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/gorilla/mux"
	"github.com/micro/cli/v2"
	ahandler "github.com/micro/go-micro/v3/api/handler"
	aapi "github.com/micro/go-micro/v3/api/handler/api"
	"github.com/micro/go-micro/v3/api/handler/event"
	ahttp "github.com/micro/go-micro/v3/api/handler/http"
	arpc "github.com/micro/go-micro/v3/api/handler/rpc"
	"github.com/micro/go-micro/v3/api/handler/web"
	"github.com/micro/go-micro/v3/api/resolver"
	"github.com/micro/go-micro/v3/api/resolver/grpc"
	"github.com/micro/go-micro/v3/api/resolver/host"
	"github.com/micro/go-micro/v3/api/resolver/path"
	"github.com/micro/go-micro/v3/api/router"
	regRouter "github.com/micro/go-micro/v3/api/router/registry"
	"github.com/micro/go-micro/v3/api/server"
	"github.com/micro/go-micro/v3/api/server/acme"
	"github.com/micro/go-micro/v3/api/server/acme/autocert"
	"github.com/micro/go-micro/v3/api/server/acme/certmagic"
	httpapi "github.com/micro/go-micro/v3/api/server/http"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/sync/memory"
	"github.com/micro/micro/v3/client"
	"github.com/micro/micro/v3/client/api/auth"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/handler"
	"github.com/micro/micro/v3/internal/helper"
	rrmicro "github.com/micro/micro/v3/internal/resolver/api"
	"github.com/micro/micro/v3/internal/stats"
	"github.com/micro/micro/v3/service"
	muregistry "github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/store"
)

var (
	Name                  = "go.micro.api"
	Address               = ":8080"
	Handler               = "meta"
	Resolver              = "micro"
	RPCPath               = "/rpc"
	APIPath               = "/"
	ProxyPath             = "/{service:[a-zA-Z0-9]+}"
	Namespace             = "go.micro"
	HeaderPrefix          = "X-Micro-"
	EnableRPC             = false
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA
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
	if len(ctx.String("enable_rpc")) > 0 {
		EnableRPC = ctx.Bool("enable_rpc")
	}
	if len(ctx.String("acme_provider")) > 0 {
		ACMEProvider = ctx.String("acme_provider")
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
	if len(ctx.String("api_namespace")) > 0 {
		Namespace = ctx.String("api_namespace")
	}
	// initialise service
	srv := service.New(service.Name(Name))

	// Init API
	var opts []server.Option

	if ctx.Bool("enable_acme") {
		hosts := helper.ACMEHosts(ctx)
		opts = append(opts, server.EnableACME(true))
		opts = append(opts, server.ACMEHosts(hosts...))
		switch ACMEProvider {
		case "autocert":
			opts = append(opts, server.ACMEProvider(autocert.NewProvider()))
		case "certmagic":
			if ACMEChallengeProvider != "cloudflare" {
				log.Fatal("The only implemented DNS challenge provider is cloudflare")
			}

			apiToken := os.Getenv("CF_API_TOKEN")
			if len(apiToken) == 0 {
				log.Fatal("env variables CF_API_TOKEN and CF_ACCOUNT_ID must be set")
			}

			storage := certmagic.NewStorage(
				memory.NewSync(),
				store.DefaultStore,
			)

			config := cloudflare.NewDefaultConfig()
			config.AuthToken = apiToken
			config.ZoneToken = apiToken
			challengeProvider, err := cloudflare.NewDNSProviderConfig(config)
			if err != nil {
				log.Fatal(err.Error())
			}

			opts = append(opts,
				server.ACMEProvider(
					certmagic.NewProvider(
						acme.AcceptToS(true),
						acme.CA(ACMECA),
						acme.Cache(storage),
						acme.ChallengeProvider(challengeProvider),
						acme.OnDemand(false),
					),
				),
			)
		default:
			log.Fatalf("%s is not a valid ACME provider\n", ACMEProvider)
		}
	} else if ctx.Bool("enable_tls") {
		config, err := helper.TLSConfig(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		opts = append(opts, server.EnableTLS(true))
		opts = append(opts, server.TLSConfig(config))
	}

	if ctx.Bool("enable_cors") {
		opts = append(opts, server.EnableCORS(true))
	}

	// create the router
	var h http.Handler
	r := mux.NewRouter()
	h = r

	if ctx.Bool("enable_stats") {
		st := stats.New()
		r.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(r)
		st.Start()
		defer st.Stop()
	}

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
	rr := rrmicro.NewResolver(ropts...)

	switch Resolver {
	case "host":
		rr = host.NewResolver(ropts...)
	case "path":
		rr = path.NewResolver(ropts...)
	case "grpc":
		rr = grpc.NewResolver(ropts...)
	}

	// register rpc handler
	if EnableRPC {
		log.Infof("Registering RPC Handler at %s", RPCPath)
		r.Handle(RPCPath, handler.NewRPCHandler(rr))
	}

	switch Handler {
	case "rpc":
		log.Infof("Registering API RPC Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(arpc.Handler),
			router.WithResolver(rr),
			router.WithRegistry(muregistry.DefaultRegistry),
		)
		rp := arpc.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(rp)
	case "api":
		log.Infof("Registering API Request Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(aapi.Handler),
			router.WithResolver(rr),
			router.WithRegistry(muregistry.DefaultRegistry),
		)
		ap := aapi.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(ap)
	case "event":
		log.Infof("Registering API Event Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(event.Handler),
			router.WithResolver(rr),
			router.WithRegistry(muregistry.DefaultRegistry),
		)
		ev := event.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(ev)
	case "http", "proxy":
		log.Infof("Registering API HTTP Handler at %s", ProxyPath)
		rt := regRouter.NewRouter(
			router.WithHandler(ahttp.Handler),
			router.WithResolver(rr),
			router.WithRegistry(muregistry.DefaultRegistry),
		)
		ht := ahttp.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithClient(srv.Client()),
		)
		r.PathPrefix(ProxyPath).Handler(ht)
	case "web":
		log.Infof("Registering API Web Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithHandler(web.Handler),
			router.WithResolver(rr),
			router.WithRegistry(muregistry.DefaultRegistry),
		)
		w := web.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithClient(srv.Client()),
		)
		r.PathPrefix(APIPath).Handler(w)
	default:
		log.Infof("Registering API Default Handler at %s", APIPath)
		rt := regRouter.NewRouter(
			router.WithResolver(rr),
			router.WithRegistry(muregistry.DefaultRegistry),
		)
		r.PathPrefix(APIPath).Handler(handler.Meta(srv, rt, Namespace))
	}

	// create the auth wrapper and the server
	authWrapper := auth.Wrapper(rr, Namespace)
	api := httpapi.NewServer(Address, server.WrapHandler(authWrapper))

	api.Init(opts...)
	api.Handle("/", h)

	// Start API
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "api",
		Usage:  "Run the api gateway",
		Action: Run,
		Flags: append(client.Flags,
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
				Name:    "type",
				Usage:   "Set the service type used by the API e.g. api",
				EnvVars: []string{"MICRO_API_TYPE"},
			},
			&cli.StringFlag{
				Name:    "resolver",
				Usage:   "Set the hostname resolver used by the API {host, path, grpc}",
				EnvVars: []string{"MICRO_API_RESOLVER"},
			},
			&cli.BoolFlag{
				Name:    "enable_rpc",
				Usage:   "Enable call the backend directly via /rpc",
				EnvVars: []string{"MICRO_API_ENABLE_RPC"},
			},
			&cli.BoolFlag{
				Name:    "enable_cors",
				Usage:   "Enable CORS, allowing the API to be called by frontend applications",
				EnvVars: []string{"MICRO_API_ENABLE_CORS"},
				Value:   true,
			},
		),
	})
}
