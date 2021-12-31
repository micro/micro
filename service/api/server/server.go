// Package api is an API Gateway
package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/gorilla/mux"
	"github.com/micro/micro/v3/plugin"
	pb "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service"
	apiserver "github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/api/auth"
	ahandler "github.com/micro/micro/v3/service/api/handler"
	aapi "github.com/micro/micro/v3/service/api/handler/api"
	"github.com/micro/micro/v3/service/api/handler/event"
	ahttp "github.com/micro/micro/v3/service/api/handler/http"
	arpc "github.com/micro/micro/v3/service/api/handler/rpc"
	"github.com/micro/micro/v3/service/api/handler/web"
	"github.com/micro/micro/v3/service/api/resolver"
	"github.com/micro/micro/v3/service/api/resolver/grpc"
	"github.com/micro/micro/v3/service/api/resolver/host"
	"github.com/micro/micro/v3/service/api/resolver/path"
	"github.com/micro/micro/v3/service/api/resolver/subdomain"
	"github.com/micro/micro/v3/service/api/router"
	regRouter "github.com/micro/micro/v3/service/api/router/registry"
	httpapi "github.com/micro/micro/v3/service/api/server/http"
	log "github.com/micro/micro/v3/service/logger"
	muregistry "github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/util/acme"
	"github.com/micro/micro/v3/util/acme/autocert"
	"github.com/micro/micro/v3/util/acme/certmagic"
	"github.com/micro/micro/v3/util/helper"
	"github.com/micro/micro/v3/util/opentelemetry"
	"github.com/micro/micro/v3/util/opentelemetry/jaeger"
	"github.com/micro/micro/v3/util/sync/memory"
	"github.com/micro/micro/v3/util/wrapper"
	"github.com/opentracing/opentracing-go"
	"github.com/urfave/cli/v2"
)

var (
	Name                  = "api"
	Address               = ":8080"
	Handler               = "meta"
	Resolver              = "micro"
	APIPath               = "/"
	ProxyPath             = "/{service:[a-zA-Z0-9]+}"
	Namespace             = ""
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA
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
		&cli.BoolFlag{
			Name:    "enable_acme",
			Usage:   "Enables ACME support via Let's Encrypt. ACME hosts should also be specified.",
			EnvVars: []string{"MICRO_API_ENABLE_ACME"},
		},
		&cli.StringFlag{
			Name:    "acme_hosts",
			Usage:   "Comma separated list of hostnames to manage ACME certs for",
			EnvVars: []string{"MICRO_API_ACME_HOSTS"},
		},
		&cli.StringFlag{
			Name:    "acme_provider",
			Usage:   "The provider that will be used to communicate with Let's Encrypt. Valid options: autocert, certmagic",
			EnvVars: []string{"MICRO_API_ACME_PROVIDER"},
		},
		&cli.BoolFlag{
			Name:    "enable_tls",
			Usage:   "Enable TLS support. Expects cert and key file to be specified",
			EnvVars: []string{"MICRO_API_ENABLE_TLS"},
		},
		&cli.StringFlag{
			Name:    "tls_cert_file",
			Usage:   "Path to the TLS Certificate file",
			EnvVars: []string{"MICRO_API_TLS_CERT_FILE"},
		},
		&cli.StringFlag{
			Name:    "tls_key_file",
			Usage:   "Path to the TLS Key file",
			EnvVars: []string{"MICRO_API_TLS_KEY_FILE"},
		},
		&cli.StringFlag{
			Name:    "tls_client_ca_file",
			Usage:   "Path to the TLS CA file to verify clients against",
			EnvVars: []string{"MICRO_API_TLS_CLIENT_CA_FILE"},
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
	// initialise service
	srv := service.New(service.Name(Name))

	// Init API
	var opts []apiserver.Option

	if ctx.Bool("enable_acme") {
		hosts := helper.ACMEHosts(ctx)
		opts = append(opts, apiserver.EnableACME(true))
		opts = append(opts, apiserver.ACMEHosts(hosts...))
		switch ACMEProvider {
		case "autocert":
			opts = append(opts, apiserver.ACMEProvider(autocert.NewProvider()))
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
				apiserver.ACMEProvider(
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

		opts = append(opts, apiserver.EnableTLS(true))
		opts = append(opts, apiserver.TLSConfig(config))
	}

	if ctx.Bool("enable_cors") {
		opts = append(opts, apiserver.EnableCORS(true))
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
	case "http":
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
		r.PathPrefix(APIPath).Handler(Meta(srv, rt, Namespace))
	}

	// register all the http handler plugins
	for _, p := range plugin.Plugins() {
		if v := p.Handler(); v != nil {
			h = v(h)
		}
	}

	reporterAddress := ctx.String("tracing_reporter_address")
	if len(reporterAddress) == 0 {
		reporterAddress = jaeger.DefaultReporterAddress
	}
	// Create a new Jaeger opentracer:
	openTracer, traceCloser, err := jaeger.New(
		opentelemetry.WithServiceName("API"),
		opentelemetry.WithTraceReporterAddress(reporterAddress),
	)
	log.Infof("Setting jaeger global tracer to %s", reporterAddress)
	defer traceCloser.Close() // Make sure we flush any pending traces before shutdown:
	if err != nil {
		log.Warnf("Unable to prepare a Jaeger tracer: %s", err)
	} else {
		// Set the global default opentracing tracer:
		opentracing.SetGlobalTracer(openTracer)
	}
	opentelemetry.DefaultOpenTracer = openTracer

	// append the opentelemetry wrapper
	h = wrapper.HTTPWrapper(h)

	// append the auth wrapper
	h = auth.Wrapper(rr, Namespace)(h)

	// create a new api server with wrappers
	api := httpapi.NewServer(Address)
	// initialise
	api.Init(opts...)
	// register the handler
	api.Handle("/", h)

	// Start API
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

	pb.RegisterApiHandler(srv.Server(), &ahandler.APIHandler{})

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
