// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/auth/jwt"
	"github.com/micro/go-micro/v3/auth/noop"
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/broker/http"
	"github.com/micro/go-micro/v3/broker/nats"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/registry/etcd"
	"github.com/micro/go-micro/v3/registry/mdns"
	"github.com/micro/go-micro/v3/registry/memory"
	"github.com/micro/go-micro/v3/router"
	regRouter "github.com/micro/go-micro/v3/router/registry"
	"github.com/micro/go-micro/v3/runtime/kubernetes"
	"github.com/micro/go-micro/v3/runtime/local"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/go-micro/v3/store/cockroach"
	"github.com/micro/go-micro/v3/store/file"
	mem "github.com/micro/go-micro/v3/store/memory"
	"github.com/micro/micro/v3/service/logger"

	inAuth "github.com/micro/micro/v3/internal/auth"
	microAuth "github.com/micro/micro/v3/service/auth"
	microBroker "github.com/micro/micro/v3/service/broker"
	microClient "github.com/micro/micro/v3/service/client"
	microConfig "github.com/micro/micro/v3/service/config"
	microRegistry "github.com/micro/micro/v3/service/registry"
	microRouter "github.com/micro/micro/v3/service/router"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	microServer "github.com/micro/micro/v3/service/server"
	microStore "github.com/micro/micro/v3/service/store"
)

// profiles which when called will configure micro to run in that environment
var profiles = map[string]*Profile{
	// built in profiles
	"ci":         CI,
	"test":       Test,
	"local":      Local,
	"kubernetes": Kubernetes,
	"platform":   Platform,
	"client":     Client,
	"service":    Service,
}

// Profile configures an environment
type Profile struct {
	// name of the profile
	Name string
	// network address to use
	NetworkAddress string
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

// CI profile to use for CI tests
var CI = &Profile{
	Name:           "ci",
	NetworkAddress: "127.0.0.1:8085",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = jwt.NewAuth()
		microRuntime.DefaultRuntime = local.NewRuntime()
		microStore.DefaultStore = file.NewStore()
		microConfig.DefaultConfig, _ = config.NewConfig()
		setBroker(http.NewBroker())
		setRegistry(etcd.NewRegistry())
		setupJWTRules()
		return nil
	},
}

// Client profile is for any entrypoint that behaves as a client
var Client = &Profile{
	Name:  "client",
	Setup: func(ctx *cli.Context) error { return nil },
}

// Local profile to run locally
var Local = &Profile{
	Name:           "local",
	NetworkAddress: "127.0.0.0:8085",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = noop.NewAuth()
		microRuntime.DefaultRuntime = local.NewRuntime()
		microStore.DefaultStore = file.NewStore()
		microConfig.DefaultConfig, _ = config.NewConfig()
		setBroker(http.NewBroker())
		setRegistry(mdns.NewRegistry())
		setupJWTRules()
		return nil
	},
}

// Kubernetes profile to run on kubernetes
var Kubernetes = &Profile{
	Name:           "kubernetes",
	NetworkAddress: "micro-network.default.svc.cluster.local:8085",
	Setup: func(ctx *cli.Context) error {
		// TODO: implement
		// registry kubernetes
		// router static
		// config configmap
		// store ...
		microAuth.DefaultAuth = jwt.NewAuth()
		setupJWTRules()
		return nil
	},
}

// Platform is for running the micro platform
var Platform = &Profile{
	Name:           "platform",
	NetworkAddress: "micro-network.default.svc.cluster.local:8085",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = jwt.NewAuth()
		microConfig.DefaultConfig, _ = config.NewConfig()
		microRuntime.DefaultRuntime = kubernetes.NewRuntime()
		setBroker(nats.NewBroker(broker.Addrs("nats-cluster")))
		setRegistry(etcd.NewRegistry(registry.Addrs("etcd-cluster")))
		setupJWTRules()

		// the cockroach store will connect immediately so the address must be passed
		// when the store is created. The cockroach store address contains the location
		// of certs so it can't be defaulted like the broker and registry.
		microStore.DefaultStore = cockroach.NewStore(store.Nodes(ctx.String("store_address")))
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
		microConfig.DefaultConfig, _ = config.NewConfig()
		setRegistry(memory.NewRegistry())
		return nil
	},
}

func setRegistry(reg registry.Registry) {
	microRegistry.DefaultRegistry = reg
	microRouter.DefaultRouter = regRouter.NewRouter()
	microRouter.DefaultRouter.Init(router.Registry(reg))
	microServer.DefaultServer.Init(server.Registry(reg))
	microClient.DefaultClient.Init(client.Registry(reg))
}

func setBroker(b broker.Broker) {
	microBroker.DefaultBroker = b
	microClient.DefaultClient.Init(client.Broker(b))
	microServer.DefaultServer.Init(server.Broker(b))
}

func setupJWTRules() {
	for _, rule := range inAuth.SystemRules {
		if err := microAuth.DefaultAuth.Grant(rule); err != nil {
			logger.Fatal("Error creating default rule: %v", err)
		}
	}
}
