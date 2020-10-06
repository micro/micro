// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"
	"os"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/auth/jwt"
	"github.com/micro/go-micro/v3/auth/noop"
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/broker/http"
	memBroker "github.com/micro/go-micro/v3/broker/memory"
	"github.com/micro/go-micro/v3/client"
	config "github.com/micro/go-micro/v3/config/store"
	memStream "github.com/micro/go-micro/v3/events/stream/memory"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/registry/mdns"
	"github.com/micro/go-micro/v3/registry/memory"
	"github.com/micro/go-micro/v3/router"
	k8sRouter "github.com/micro/go-micro/v3/router/kubernetes"
	regRouter "github.com/micro/go-micro/v3/router/registry"
	"github.com/micro/go-micro/v3/runtime/kubernetes"
	"github.com/micro/go-micro/v3/runtime/local"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/store/file"
	mem "github.com/micro/go-micro/v3/store/memory"
	"github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"

	inAuth "github.com/micro/micro/v3/internal/auth"
	"github.com/micro/micro/v3/internal/user"
	microAuth "github.com/micro/micro/v3/service/auth"
	microBroker "github.com/micro/micro/v3/service/broker"
	microClient "github.com/micro/micro/v3/service/client"
	microConfig "github.com/micro/micro/v3/service/config"
	microEvents "github.com/micro/micro/v3/service/events"
	microRegistry "github.com/micro/micro/v3/service/registry"
	microRouter "github.com/micro/micro/v3/service/router"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	microServer "github.com/micro/micro/v3/service/server"
	microStore "github.com/micro/micro/v3/service/store"
)

// profiles which when called will configure micro to run in that environment
var profiles = map[string]*Profile{
	// built in profiles
	"client":     Client,
	"service":    Service,
	"test":       Test,
	"local":      Local,
	"kubernetes": Kubernetes,
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

// Local profile to run locally
var Local = &Profile{
	Name: "local",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = jwt.NewAuth()
		microStore.DefaultStore = file.NewStore()
		SetupConfigSecretKey(ctx)
		microConfig.DefaultConfig, _ = config.NewConfig(microStore.DefaultStore, "")
		SetupBroker(http.NewBroker())
		SetupRouter(regRouter.NewRouter())
		SetupRegistry(mdns.NewRegistry())
		SetupJWT(ctx)

		// use the local runtime, note: the local runtime is designed to run source code directly so
		// the runtime builder should NOT be set when using this implementation
		microRuntime.DefaultRuntime = local.NewRuntime()

		var err error
		microEvents.DefaultStream, err = memStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}

		microStore.DefaultBlobStore, err = file.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		return nil
	},
}

// Kubernetes profile to run on kubernetes with zero deps. Designed for use with the micro helm chart
var Kubernetes = &Profile{
	Name: "kubernetes",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = jwt.NewAuth()
		SetupJWT(ctx)

		microRuntime.DefaultRuntime = kubernetes.NewRuntime()
		SetupRouter(k8sRouter.NewRouter())

		var err error
		microEvents.DefaultStream, err = memStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}

		microStore.DefaultStore = file.NewStore(file.WithDir("/store"))
		microStore.DefaultBlobStore, err = file.NewBlobStore(file.WithDir("/store/blob"))
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		// the registry service uses the memory registry, the other core services will use the default
		// rpc client and call the registry service
		if ctx.Args().Get(1) == "registry" {
			SetupRegistry(memory.NewRegistry())
		}

		// the broker service uses the memory broker, the other core services will use the default
		// rpc client and call the broker service
		if ctx.Args().Get(1) == "broker" {
			SetupBroker(memBroker.NewBroker())
		}

		microConfig.DefaultConfig, err = config.NewConfig(microStore.DefaultStore, "")
		if err != nil {
			logger.Fatalf("Error configuring config: %v", err)
		}
		SetupConfigSecretKey(ctx)

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
		microConfig.DefaultConfig, _ = config.NewConfig(microStore.DefaultStore, "")
		SetupRouter(regRouter.NewRouter())
		SetupRegistry(memory.NewRegistry())
		return nil
	},
}

// SetupRegistry configures the registry
func SetupRegistry(reg registry.Registry) {
	microRegistry.DefaultRegistry = reg
	microRouter.DefaultRouter.Init(router.Registry(reg))
	microServer.DefaultServer.Init(server.Registry(reg))
	microClient.DefaultClient.Init(client.Registry(reg))
}

// SetupRouter configures the router. Should be called before SetupRegistry.
func SetupRouter(rtr router.Router) {
	microRouter.DefaultRouter = rtr
	microClient.DefaultClient.Init(client.Router(rtr))
}

// SetupBroker configures the broker
func SetupBroker(b broker.Broker) {
	microBroker.DefaultBroker = b
	microClient.DefaultClient.Init(client.Broker(b))
	microServer.DefaultServer.Init(server.Broker(b))
}

// SetupJWT configures the default internal system rules
func SetupJWT(ctx *cli.Context) {
	for _, rule := range inAuth.SystemRules {
		if err := microAuth.DefaultAuth.Grant(rule); err != nil {
			logger.Fatal("Error creating default rule: %v", err)
		}
	}

	// Generate public and private keys if none are provided
	pubKey := ctx.String("auth_public_key")
	privKey := ctx.String("auth_private_key")
	if len(privKey) == 0 || len(pubKey) == 0 {
		var err error
		privKey, pubKey, err = user.GetJWTCerts()
		if err != nil {
			logger.Fatalf("Error getting keys: %v", err)
		}
	}
	microAuth.DefaultAuth.Init(
		auth.PrivateKey(privKey),
		auth.PublicKey(pubKey),
	)
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
