// Package platform is a profile for running a highly available Micro platform
package platform

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"

	"github.com/micro/micro/v3/service/runtime/kubernetes"
	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/auth/jwt"
	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/config"
	storeConfig "github.com/micro/micro/v3/service/config/store"
	"github.com/micro/micro/v3/service/events"
	evStore "github.com/micro/micro/v3/service/events/store"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/metrics"
	"github.com/micro/micro/v3/service/registry"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	microBuilder "github.com/micro/micro/v3/service/runtime/builder"
	buildSrv "github.com/micro/micro/v3/service/runtime/builder/client"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/s3"
	"github.com/urfave/cli/v2"

	// plugins
	"github.com/micro/micro/plugin/cockroach/v3"
	"github.com/micro/micro/plugin/etcd/v3"
	natsBroker "github.com/micro/micro/plugin/nats/broker/v3"
	natsStream "github.com/micro/micro/plugin/nats/stream/v3"
	"github.com/micro/micro/plugin/prometheus/v3"
)

func init() {
	profile.Register("platform", Profile)
}

// Profile is for running the micro platform
var Profile = &profile.Profile{
	Name: "platform",
	Setup: func(ctx *cli.Context) error {
		auth.DefaultAuth = jwt.NewAuth()
		// the cockroach store will connect immediately so the address must be passed
		// when the store is created. The cockroach store address contains the location
		// of certs so it can't be defaulted like the broker and registry.
		store.DefaultStore = cockroach.NewStore(store.Nodes(ctx.String("store_address")))
		config.DefaultConfig, _ = storeConfig.NewConfig(store.DefaultStore, "")
		profile.SetupBroker(natsBroker.NewBroker(broker.Addrs("nats-cluster")))
		profile.SetupRegistry(etcd.NewRegistry(registry.Addrs("etcd-cluster")))
		profile.SetupJWT(ctx)
		profile.SetupConfigSecretKey(ctx)

		// Set up a default metrics reporter (being careful not to clash with any that have already been set):
		if !metrics.IsSet() {
			prometheusReporter, err := prometheus.New()
			if err != nil {
				return err
			}
			metrics.SetDefaultMetricsReporter(prometheusReporter)
		}

		var err error
		events.DefaultStream, err = natsStream.NewStream(natsStreamOpts(ctx)...)
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}

		// only configure the blob store for the store and runtime services
		if ctx.Args().Get(1) == "runtime" || ctx.Args().Get(1) == "store" {
			store.DefaultBlobStore, err = s3.NewBlobStore(
				s3.Credentials(
					os.Getenv("MICRO_BLOB_STORE_ACCESS_KEY"),
					os.Getenv("MICRO_BLOB_STORE_SECRET_KEY"),
				),
				s3.Endpoint("minio-cluster:9000"),
				s3.Region(os.Getenv("MICRO_BLOB_STORE_REGION")),
				s3.Insecure(),
			)
			if err != nil {
				logger.Fatalf("Error configuring s3 blob store: %v", err)
			}
		}

		microBuilder.DefaultBuilder = buildSrv.NewBuilder()
		microRuntime.DefaultRuntime = kubernetes.NewRuntime()
		events.DefaultStore = evStore.NewStore(evStore.WithStore(store.DefaultStore))
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
