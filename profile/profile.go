// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"

	"github.com/micro/go-micro/v3/auth/jwt"
	"github.com/micro/go-micro/v3/auth/noop"
	"github.com/micro/go-micro/v3/broker/http"
	"github.com/micro/go-micro/v3/broker/nats"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/registry/etcd"
	"github.com/micro/go-micro/v3/registry/mdns"
	"github.com/micro/go-micro/v3/registry/memory"
	"github.com/micro/go-micro/v3/router"
	regRouter "github.com/micro/go-micro/v3/router/registry"
	"github.com/micro/go-micro/v3/runtime/kubernetes"
	"github.com/micro/go-micro/v3/runtime/local"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/store/cockroach"
	"github.com/micro/go-micro/v3/store/file"
	mem "github.com/micro/go-micro/v3/store/memory"

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
var profiles = map[string]Profile{
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
type Profile func() error

// Register a profile
func Register(name string, p Profile) error {
	if _, ok := profiles[name]; ok {
		return fmt.Errorf("profile %s already exists", name)
	}
	profiles[name] = p
	return nil
}

// Load a profile
func Load(name string) (Profile, error) {
	v, ok := profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %s does not exist", name)
	}
	return v, nil
}

// CI profile to use for CI tests
var CI Profile = func() error {
	microAuth.DefaultAuth = jwt.NewAuth()
	microBroker.DefaultBroker = http.NewBroker()
	microRuntime.DefaultRuntime = local.NewRuntime()
	microStore.DefaultStore = file.NewStore()
	microConfig.DefaultConfig, _ = config.NewConfig()
	setRegistry(etcd.NewRegistry())
	setupJWTRules()
	return nil
}

// Client profile is for any entrypoint that behaves as a client
var Client Profile = func() error {
	// Defaults to service implementations
	return nil
}

// Local profile to run locally
var Local Profile = func() error {
	microAuth.DefaultAuth = noop.NewAuth()
	microBroker.DefaultBroker = http.NewBroker()
	microRuntime.DefaultRuntime = local.NewRuntime()
	microStore.DefaultStore = file.NewStore()
	microConfig.DefaultConfig, _ = config.NewConfig()
	setRegistry(mdns.NewRegistry())
	setupJWTRules()
	return nil
}

// Kubernetes profile to run on kubernetes
var Kubernetes Profile = func() error {
	// TODO: implement
	// auth jwt
	// registry kubernetes
	// router static
	// config configmap
	// store ...
	microAuth.DefaultAuth = jwt.NewAuth()
	setupJWTRules()
	return nil
}

// Platform is for running the micro platform
var Platform Profile = func() error {
	microAuth.DefaultAuth = jwt.NewAuth()
	microBroker.DefaultBroker = nats.NewBroker()
	microRuntime.DefaultRuntime = kubernetes.NewRuntime()
	microStore.DefaultStore = cockroach.NewStore()
	microConfig.DefaultConfig, _ = config.NewConfig()
	setRegistry(etcd.NewRegistry())
	setupJWTRules()
	return nil
}

// Service is the default for any services run
var Service Profile = func() error {
	// All values are set by default
	// Potentially better set here
	// Add any other initialisation necessary
	return nil
}

// Test profile is used for the go test suite
var Test Profile = func() error {
	microAuth.DefaultAuth = noop.NewAuth()
	microStore.DefaultStore = mem.NewStore()
	microConfig.DefaultConfig, _ = config.NewConfig()
	setRegistry(memory.NewRegistry())
	return nil
}

func setRegistry(reg registry.Registry) {
	microRegistry.DefaultRegistry = reg
	microRouter.DefaultRouter = regRouter.NewRouter()
	microRouter.DefaultRouter.Init(router.Registry(reg))
	microServer.DefaultServer.Init(server.Registry(reg))
	microClient.DefaultClient.Init(client.Registry(reg))
}

func setupJWTRules() {
	for _, rule := range inAuth.SystemRules {
		if err := microAuth.DefaultAuth.Grant(rule); err != nil {
			logger.Fatal("Error creating default rule: %v", err)
		}
	}
}
