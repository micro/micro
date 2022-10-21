// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/micro/micro/v3/service/auth/jwt"
	"github.com/micro/micro/v3/service/auth/noop"
	"github.com/micro/micro/v3/service/broker"
	memBroker "github.com/micro/micro/v3/service/broker/memory"
	"github.com/micro/micro/v3/service/client"
	grpcClient "github.com/micro/micro/v3/service/client/grpc"
	"github.com/micro/micro/v3/service/config"
	storeConfig "github.com/micro/micro/v3/service/config/store"
	evStore "github.com/micro/micro/v3/service/events/store"
	memStream "github.com/micro/micro/v3/service/events/stream/memory"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/registry/memory"
	"github.com/micro/micro/v3/service/router"
	regRouter "github.com/micro/micro/v3/service/router/registry"
	"github.com/micro/micro/v3/service/runtime/local"
	"github.com/micro/micro/v3/service/server"
	grpcServer "github.com/micro/micro/v3/service/server/grpc"
	"github.com/micro/micro/v3/service/store/file"
	mem "github.com/micro/micro/v3/service/store/memory"
	"github.com/micro/micro/v3/util/opentelemetry"
	"github.com/micro/micro/v3/util/opentelemetry/jaeger"
	"github.com/urfave/cli/v2"

	microAuth "github.com/micro/micro/v3/service/auth"
	microEvents "github.com/micro/micro/v3/service/events"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	microStore "github.com/micro/micro/v3/service/store"
	inAuth "github.com/micro/micro/v3/util/auth"
	"github.com/micro/micro/v3/util/user"
)

// profiles which when called will configure micro to run in that environment
var profiles = map[string]*Profile{
	// built in profiles
	"client":  Client,
	"service": Service,
	"server":  Server,
	"test":    Test,
	"local":   Local,
}

// Profile configures an environment
type Profile struct {
	// name of the profile
	Name string
	// function used for setup
	Setup func(*cli.Context) error
	// TODO: presetup dependencies
	// e.g start resources
}

// Register a profile
func Register(name string, p *Profile) error {
	if _, ok := profiles[name]; ok {
		return fmt.Errorf("profile %s already exists", name)
	}
	profiles[name] = p
	return nil
}

// Load a profile
func Load(name string) (*Profile, error) {
	v, ok := profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %s does not exist", name)
	}
	return v, nil
}

// Client profile is for any entrypoint that behaves as a client
var Client = &Profile{
	Name:  "client",
	Setup: func(ctx *cli.Context) error { return nil },
}

// Local profile to run as a single process
var Local = &Profile{
	Name: "local",
	Setup: func(ctx *cli.Context) error {
		// set client/server
		client.DefaultClient = grpcClient.NewClient()
		server.DefaultServer = grpcServer.NewServer()

		microAuth.DefaultAuth = jwt.NewAuth()
		microStore.DefaultStore = file.NewStore(file.WithDir(filepath.Join(user.Dir, "server", "store")))
		SetupConfigSecretKey(ctx)
		config.DefaultConfig, _ = storeConfig.NewConfig(microStore.DefaultStore, "")

		SetupJWT(ctx)
		SetupRegistry(memory.NewRegistry())
		SetupBroker(memBroker.NewBroker())

		// set the store in the model
		model.DefaultModel = model.NewModel(
			model.WithStore(microStore.DefaultStore),
		)

		microRuntime.DefaultRuntime = local.NewRuntime()

		var err error
		microEvents.DefaultStream, err = memStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}
		microEvents.DefaultStore = evStore.NewStore(
			evStore.WithStore(microStore.DefaultStore),
		)

		microStore.DefaultBlobStore, err = file.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		return nil
	},
}

var Server = &Profile{
	Name: "server",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = jwt.NewAuth()
		microStore.DefaultStore = file.NewStore(file.WithDir(filepath.Join(user.Dir, "server", "store")))
		SetupConfigSecretKey(ctx)
		config.DefaultConfig, _ = storeConfig.NewConfig(microStore.DefaultStore, "")
		SetupJWT(ctx)

		// the registry service uses the memory registry, the other core services will use the default
		// rpc client and call the registry service
		if ctx.Args().Get(1) == "registry" {
			SetupRegistry(memory.NewRegistry())
		} else {
			// set the registry address
			registry.DefaultRegistry.Init(
				registry.Addrs("localhost:8000"),
			)

			SetupRegistry(registry.DefaultRegistry)
		}

		// the broker service uses the memory broker, the other core services will use the default
		// rpc client and call the broker service
		if ctx.Args().Get(1) == "broker" {
			SetupBroker(memBroker.NewBroker())
		} else {
			broker.DefaultBroker.Init(
				broker.Addrs("localhost:8003"),
			)
			SetupBroker(broker.DefaultBroker)
		}

		// set the store in the model
		model.DefaultModel = model.NewModel(
			model.WithStore(microStore.DefaultStore),
		)

		// use the local runtime, note: the local runtime is designed to run source code directly so
		// the runtime builder should NOT be set when using this implementation
		microRuntime.DefaultRuntime = local.NewRuntime()

		var err error
		microEvents.DefaultStream, err = memStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}
		microEvents.DefaultStore = evStore.NewStore(
			evStore.WithStore(microStore.DefaultStore),
		)

		microStore.DefaultBlobStore, err = file.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		// Configure tracing with Jaeger (forced tracing):
		tracingServiceName := ctx.Args().Get(1)
		if len(tracingServiceName) == 0 {
			tracingServiceName = "Micro"
		}
		openTracer, _, err := jaeger.New(
			opentelemetry.WithServiceName(tracingServiceName),
			opentelemetry.WithSamplingRate(1),
		)
		if err != nil {
			logger.Fatalf("Error configuring opentracing: %v", err)
		}
		opentelemetry.DefaultOpenTracer = openTracer

		return nil
	},
}

// Service is the default for any services run
var Service = &Profile{
	Name:  "service",
	Setup: func(ctx *cli.Context) error { return nil },
}

// Test profile is used for the go test suite
var Test = &Profile{
	Name: "test",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = noop.NewAuth()
		microStore.DefaultStore = mem.NewStore()
		microStore.DefaultBlobStore, _ = file.NewBlobStore()
		config.DefaultConfig, _ = storeConfig.NewConfig(microStore.DefaultStore, "")
		SetupRegistry(memory.NewRegistry())
		// set the store in the model
		model.DefaultModel = model.NewModel(
			model.WithStore(microStore.DefaultStore),
		)
		return nil
	},
}

// SetupRegistry configures the registry
func SetupRegistry(reg registry.Registry) {
	registry.DefaultRegistry = reg
	router.DefaultRouter = regRouter.NewRouter(router.Registry(reg))
	client.DefaultClient.Init(client.Registry(reg), client.Router(router.DefaultRouter))
	server.DefaultServer.Init(server.Registry(reg))
}

// SetupBroker configures the broker
func SetupBroker(b broker.Broker) {
	broker.DefaultBroker = b
	client.DefaultClient.Init(client.Broker(b))
	server.DefaultServer.Init(server.Broker(b))
}

// SetupJWT configures the default internal system rules
func SetupJWT(ctx *cli.Context) {
	for _, rule := range inAuth.SystemRules {
		if err := microAuth.DefaultAuth.Grant(rule); err != nil {
			logger.Fatal("Error creating default rule: %v", err)
		}
	}
}

func SetupConfigSecretKey(ctx *cli.Context) {
	key := ctx.String("config_secret_key")
	if len(key) == 0 {
		k, err := user.GetConfigSecretKey()
		if err != nil {
			logger.Fatal("Error getting config secret: %v", err)
		}
		os.Setenv("MICRO_CONFIG_SECRET_KEY", k)
	}
}
