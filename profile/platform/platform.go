// Package platform is a profile for running a highly available Micro platform
package platform

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/micro/go-micro/v3/auth/jwt"
	"github.com/micro/go-micro/v3/broker"
	config "github.com/micro/go-micro/v3/config/store"
	evStore "github.com/micro/go-micro/v3/events/store"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/runtime/kubernetes"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"

	microAuth "github.com/micro/micro/v3/service/auth"
	microConfig "github.com/micro/micro/v3/service/config"
	microEvents "github.com/micro/micro/v3/service/events"
	microMetrics "github.com/micro/micro/v3/service/metrics"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	microStore "github.com/micro/micro/v3/service/store"

	// plugins
	"github.com/micro/go-plugins/broker/nats/v3"
	natsStream "github.com/micro/go-plugins/events/stream/nats/v3"
	metricsPrometheus "github.com/micro/go-plugins/metrics/prometheus/v3"
	"github.com/micro/go-plugins/registry/etcd/v3"
	"github.com/micro/go-plugins/store/cockroach/v3"
)

func init() {
	profile.Register("platform", Profile)
}

// Profile is for running the micro platform
var Profile = &profile.Profile{
	Name: "platform",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = jwt.NewAuth()
		// the cockroach store will connect immediately so the address must be passed
		// when the store is created. The cockroach store address contains the location
		// of certs so it can't be defaulted like the broker and registry.
		microStore.DefaultStore = cockroach.NewStore(store.Nodes(ctx.String("store_address")))
		microConfig.DefaultConfig, _ = config.NewConfig(microStore.DefaultStore, "")
		microRuntime.DefaultRuntime = kubernetes.NewRuntime()
		profile.SetupBroker(nats.NewBroker(broker.Addrs("nats-cluster")))
		profile.SetupRegistry(etcd.NewRegistry(registry.Addrs("etcd-cluster")))
		profile.SetupJWT(ctx)
		profile.SetupConfigSecretKey(ctx)

		// Set up a default metrics reporter (being careful not to clash with any that have already been set):
		if !microMetrics.IsSet() {
			prometheusReporter, err := metricsPrometheus.New()
			if err != nil {
				return err
			}
			microMetrics.SetDefaultMetricsReporter(prometheusReporter)
		}

		var err error
		microEvents.DefaultStream, err = natsStream.NewStream(natsStreamOpts(ctx)...)
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}

		microEvents.DefaultStore = evStore.NewStore(evStore.WithStore(microStore.DefaultStore))
		return nil
	},
}

// natsStreamOpts returns a slice of options which should be used to configure nats
func natsStreamOpts(ctx *cli.Context) []natsStream.Option {
	opts := []natsStream.Option{
		natsStream.Address("nats://nats-cluster:4222"),
		natsStream.ClusterID("nats-streaming-cluster"),
	}

	// Parse event TLS certs
	if len(ctx.String("events_tls_cert")) > 0 || len(ctx.String("events_tls_key")) > 0 {
		cert, err := tls.LoadX509KeyPair(ctx.String("events_tls_cert"), ctx.String("events_tls_key"))
		if err != nil {
			logger.Fatalf("Error loading event TLS cert: %v", err)
		}

		// load custom certificate authority
		caCertPool := x509.NewCertPool()
		if len(ctx.String("events_tls_ca")) > 0 {
			crt, err := ioutil.ReadFile(ctx.String("events_tls_ca"))
			if err != nil {
				logger.Fatalf("Error loading event TLS certificate authority: %v", err)
			}
			caCertPool.AppendCertsFromPEM(crt)
		}

		cfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: caCertPool}
		opts = append(opts, natsStream.TLSConfig(cfg))
	}

	return opts
}
