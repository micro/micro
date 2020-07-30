// Package proxy is a cli proxy
package proxy

import (
	"os"
	"strings"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/api/server/acme"
	"github.com/micro/go-micro/v3/api/server/acme/autocert"
	"github.com/micro/go-micro/v3/api/server/acme/certmagic"
	bmem "github.com/micro/go-micro/v3/broker/memory"
	grpcCli "github.com/micro/go-micro/v3/client/grpc"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/go-micro/v3/proxy"
	"github.com/micro/go-micro/v3/proxy/http"
	"github.com/micro/go-micro/v3/proxy/mucp"
	rmem "github.com/micro/go-micro/v3/registry/memory"
	"github.com/micro/go-micro/v3/router"
	"github.com/micro/go-micro/v3/router/registry"
	"github.com/micro/go-micro/v3/server"
	grpcSrv "github.com/micro/go-micro/v3/server/grpc"
	sgrpc "github.com/micro/go-micro/v3/server/grpc"
	"github.com/micro/go-micro/v3/sync/memory"
	"github.com/micro/go-micro/v3/util/mux"
	"github.com/micro/micro/v3/client"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/helper"
	"github.com/micro/micro/v3/service"
	muregistry "github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/store"
)

var (
	// Name of the proxy
	Name = "go.micro.proxy"
	// The address of the proxy
	Address = ":8081"
	// the proxy protocol
	Protocol = "grpc"
	// The endpoint host to route to
	Endpoint string
	// ACME (Cert management)
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA
)

func Run(ctx *cli.Context) error {
	// because MICRO_PROXY_ADDRESS is used internally by the go-micro/client
	// we need to unset it so we don't end up calling ourselves infinitely
	os.Unsetenv("MICRO_PROXY_ADDRESS")
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("proxy_address")) > 0 {
		Address = ctx.String("proxy_address")
	}
	if len(ctx.String("endpoint")) > 0 {
		Endpoint = ctx.String("endpoint")
	}
	if len(ctx.String("protocol")) > 0 {
		Protocol = ctx.String("protocol")
	}
	if len(ctx.String("acme_provider")) > 0 {
		ACMEProvider = ctx.String("acme_provider")
	}

	// new service
	service := service.New(service.Name(Name))

	// set the context
	popts := []proxy.Option{
		proxy.WithRouter(registry.NewRouter(
			router.Registry(muregistry.DefaultRegistry),
		)),
	}

	// new proxy
	var p proxy.Proxy
	// setup the default server
	var srv server.Server

	// set endpoint
	if len(Endpoint) > 0 {
		switch {
		case strings.HasPrefix(Endpoint, "grpc://"):
			ep := strings.TrimPrefix(Endpoint, "grpc://")
			popts = append(popts, proxy.WithEndpoint(ep))
			Protocol = "grpc"
		case strings.HasPrefix(Endpoint, "http://"):
			// TODO: strip prefix?
			popts = append(popts, proxy.WithEndpoint(Endpoint))
			Protocol = "http"
		default:
			// TODO: strip prefix?
			popts = append(popts, proxy.WithEndpoint(Endpoint))
			Protocol = "mucp"
		}
	}

	serverOpts := []server.Option{
		server.Address(Address),
		server.Registry(rmem.NewRegistry()),
		server.Broker(bmem.NewBroker()),
	}

	// enable acme will create a net.Listener which
	if ctx.Bool("enable_acme") {
		var ap acme.Provider

		switch ACMEProvider {
		case "autocert":
			ap = autocert.NewProvider()
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

			// define the provider
			ap = certmagic.NewProvider(
				acme.AcceptToS(true),
				acme.CA(ACMECA),
				acme.Cache(storage),
				acme.ChallengeProvider(challengeProvider),
				acme.OnDemand(false),
			)
		default:
			log.Fatalf("Unsupported acme provider: %s\n", ACMEProvider)
		}

		// generate the tls config
		config, err := ap.TLSConfig(helper.ACMEHosts(ctx)...)
		if err != nil {
			log.Fatalf("Failed to generate acme tls config: %v", err)
		}

		// set the tls config
		serverOpts = append(serverOpts, server.TLSConfig(config))
		// enable tls will leverage tls certs and generate a tls.Config
	} else if ctx.Bool("enable_tls") {
		// get certificates from the context
		config, err := helper.TLSConfig(ctx)
		if err != nil {
			log.Fatal(err)
			return err
		}
		serverOpts = append(serverOpts, server.TLSConfig(config))
	}

	// wrap the proxy using the proxy's authHandler
	authOpt := server.WrapHandler(authHandler())
	serverOpts = append(serverOpts, authOpt)

	// set proxy
	switch Protocol {
	case "http":
		p = http.NewProxy(popts...)
		serverOpts = append(serverOpts, server.WithRouter(p))
		// TODO: http server
		srv = grpcSrv.NewServer(serverOpts...)
	case "mucp":
		popts = append(popts, proxy.WithClient(grpcCli.NewClient()))
		p = mucp.NewProxy(popts...)

		serverOpts = append(serverOpts, server.WithRouter(p))
		srv = grpcSrv.NewServer(serverOpts...)
	default:
		p = mucp.NewProxy(popts...)

		serverOpts = append(serverOpts, server.WithRouter(p))
		srv = sgrpc.NewServer(serverOpts...)
	}

	if len(Endpoint) > 0 {
		log.Infof("Proxy [%s] serving endpoint: %s", p.String(), Endpoint)
	} else {
		log.Infof("Proxy [%s] serving protocol: %s", p.String(), Protocol)
	}

	// create a new proxy muxer which includes the debug handler
	muxer := mux.New(Name, p)

	// set the router
	service.Server().Init(
		server.WithRouter(muxer),
	)

	// Start the proxy server
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	// Run internal service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop the server
	if err := srv.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:  "proxy",
		Usage: "Run the service proxy",
		Flags: append(client.Flags,
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the proxy http address e.g 0.0.0.0:8081",
				EnvVars: []string{"MICRO_PROXY_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "protocol",
				Usage:   "Set the protocol used for proxying e.g mucp, grpc, http",
				EnvVars: []string{"MICRO_PROXY_PROTOCOL"},
			},
			&cli.StringFlag{
				Name:    "endpoint",
				Usage:   "Set the endpoint to route to e.g greeter or localhost:9090",
				EnvVars: []string{"MICRO_PROXY_ENDPOINT"},
			},
		),
		Action: Run,
	})
}
