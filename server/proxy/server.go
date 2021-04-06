// Package proxy is a cli proxy
package proxy

import (
	"os"
	"strings"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/micro/micro/v3/client"
	"github.com/micro/micro/v3/internal/api/server/acme"
	"github.com/micro/micro/v3/internal/api/server/acme/autocert"
	"github.com/micro/micro/v3/internal/api/server/acme/certmagic"
	"github.com/micro/micro/v3/internal/helper"
	"github.com/micro/micro/v3/internal/muxer"
	"github.com/micro/micro/v3/internal/sync/memory"
	"github.com/micro/micro/v3/service"
	bmem "github.com/micro/micro/v3/service/broker/memory"
	muclient "github.com/micro/micro/v3/service/client"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/proxy"
	"github.com/micro/micro/v3/service/proxy/grpc"
	"github.com/micro/micro/v3/service/proxy/http"
	"github.com/micro/micro/v3/service/proxy/mucp"
	"github.com/micro/micro/v3/service/registry/noop"
	murouter "github.com/micro/micro/v3/service/router"
	"github.com/micro/micro/v3/service/server"
	sgrpc "github.com/micro/micro/v3/service/server/grpc"
	"github.com/micro/micro/v3/service/store"
	"github.com/urfave/cli/v2"
)

var (
	// Name of the proxy
	Name = "proxy"
	// The address of the proxy
	Address = ":8081"
	// Is gRPCWeb enabled
	GRPCWebEnabled = false
	// The address of the proxy
	GRPCWebAddress = ":8082"
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
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if ctx.Bool("grpc-web") {
		GRPCWebEnabled = ctx.Bool("grpcWeb")
	}
	if len(ctx.String("grpc-web-port")) > 0 {
		GRPCWebAddress = ctx.String("grpcWebAddr")
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
		proxy.WithRouter(murouter.DefaultRouter),
		proxy.WithClient(muclient.DefaultClient),
	}

	// set endpoint
	if len(Endpoint) > 0 {
		ep := Endpoint

		switch {
		case strings.HasPrefix(Endpoint, "grpc://"):
			ep = strings.TrimPrefix(Endpoint, "grpc://")
			Protocol = "grpc"
		case strings.HasPrefix(Endpoint, "http://"):
			Protocol = "http"
		case strings.HasPrefix(Endpoint, "mucp://"):
			ep = strings.TrimPrefix(Endpoint, "mucp://")
			Protocol = "mucp"
		}

		popts = append(popts, proxy.WithEndpoint(ep))
	}

	serverOpts := []server.Option{
		server.Name(Name),
		server.Address(Address),
		server.Registry(noop.NewRegistry()),
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

	// new proxy
	var p proxy.Proxy

	// set proxy
	switch Protocol {
	case "http":
		p = http.NewProxy(popts...)
		// TODO: http server
	case "mucp":
		p = mucp.NewProxy(popts...)
	default:
		// default to the grpc proxy
		p = grpc.NewProxy(popts...)
	}

	// wrap the proxy using the proxy's authHandler
	authOpt := server.WrapHandler(authHandler())
	serverOpts = append(serverOpts, authOpt)
	serverOpts = append(serverOpts, server.WithRouter(p))

	if len(Endpoint) > 0 {
		log.Infof("Proxy [%s] serving endpoint: %s", p.String(), Endpoint)
	} else {
		log.Infof("Proxy [%s] serving protocol: %s", p.String(), Protocol)
	}

	if GRPCWebEnabled {
		serverOpts = append(serverOpts, sgrpc.GRPCWebPort(GRPCWebAddress))
		serverOpts = append(serverOpts, sgrpc.GRPCWebOptions(
			grpcweb.WithCorsForRegisteredEndpointsOnly(false),
			grpcweb.WithOriginFunc(func(origin string) bool { return true })))

		log.Infof("Proxy [%s] serving gRPC-Web on %s", p.String(), GRPCWebAddress)
	}

	// create a new grpc server
	srv := sgrpc.NewServer(serverOpts...)

	// create a new proxy muxer which includes the debug handler
	muxer := muxer.New(Name, p)

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

var (
	Flags = append(client.Flags,
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
		&cli.BoolFlag{
			Name:    "grpc-web",
			Usage:   "Enable the gRPCWeb server",
			EnvVars: []string{"MICRO_PROXY_GRPC_WEB"},
		},
		&cli.StringFlag{
			Name:    "grpc-web-addr",
			Usage:   "Set the gRPC web addr on the proxy",
			EnvVars: []string{"MICRO_PROXY_GRPC_WEB_ADDRESS"},
		},
	)
)
