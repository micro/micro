// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/auth/jwt"
	"github.com/micro/go-micro/v3/auth/noop"
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/broker/http"
	"github.com/micro/go-micro/v3/broker/nats"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/registry/etcd"
	"github.com/micro/go-micro/v3/registry/mdns"
	memReg "github.com/micro/go-micro/v3/registry/memory"
	"github.com/micro/go-micro/v3/router"
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/runtime/kubernetes"
	"github.com/micro/go-micro/v3/runtime/local"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/go-micro/v3/store/cockroach"
	"github.com/micro/go-micro/v3/store/file"
	mem "github.com/micro/go-micro/v3/store/memory"

	microAuth "github.com/micro/micro/v3/service/auth"
	microRouter "github.com/micro/micro/v3/service/router"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	microStore "github.com/micro/micro/v3/service/store"

	authClient "github.com/micro/micro/v3/service/auth/client"
	brokerClient "github.com/micro/micro/v3/service/broker/client"
	registryClient "github.com/micro/micro/v3/service/registry/client"
	routerClient "github.com/micro/micro/v3/service/router/client"
	runtimeClient "github.com/micro/micro/v3/service/runtime/client"
	storeClient "github.com/micro/micro/v3/service/store/client"
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

// Profile configures an environment. If an implementation is
// not specified, the RPC implementation will be used
type Profile struct {
	Auth       func(opts ...auth.Option) auth.Auth
	Broker     func(opts ...broker.Option) broker.Broker
	Registry   func(opts ...registry.Option) registry.Registry
	Router     func(opts ...router.Option) router.Router
	Runtime    func(opts ...runtime.Option) runtime.Runtime
	Store      func(opts ...store.Option) store.Store
	AfterSetup func()
}

// Setup the profile
func (p *Profile) Setup() {
	if p.Auth == nil {
		microAuth.DefaultAuth = authClient.NewAuth()
	} else {
		microAuth.DefaultAuth = p.Auth()
	}

	if p.Broker == nil {
		setBroker(brokerClient.NewBroker())
	} else {
		setBroker(p.Broker())
	}

	if p.Registry == nil {
		setRegistry(registryClient.NewRegistry())
	} else {
		setRegistry(p.Registry())
	}

	if p.Router == nil {
		microRouter.DefaultRouter = routerClient.NewRouter()
	} else {
		microRouter.DefaultRouter = p.Router()
	}

	if p.Runtime == nil {
		microRuntime.DefaultRuntime = runtimeClient.NewRuntime()
	} else {
		microRuntime.DefaultRuntime = p.Runtime()
	}

	if p.Store == nil {
		microStore.DefaultStore = storeClient.NewStore()
	} else {
		microStore.DefaultStore = p.Store()
	}

	if p.AfterSetup != nil {
		p.AfterSetup()
	}
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
	Auth:       jwt.NewAuth,
	Runtime:    local.NewRuntime,
	Store:      file.NewStore,
	Broker:     http.NewBroker,
	Registry:   etcd.NewRegistry,
	AfterSetup: setupJWTRules,
}

// Client profile is for any entrypoint that behaves as a client
var Client = &Profile{}

// Local profile to run locally
var Local = &Profile{
	Auth:     noop.NewAuth,
	Broker:   http.NewBroker,
	Store:    file.NewStore,
	Runtime:  local.NewRuntime,
	Registry: mdns.NewRegistry,
}

// Kubernetes profile to run on kubernetes
var Kubernetes = &Profile{
	Auth: jwt.NewAuth,
	// TODO: implement
	// registry kubernetes
	// router static
	// config configmap
	// store ...
	AfterSetup: setupJWTRules,
}

// Platform is for running the micro platform
var Platform = &Profile{
	Auth:    jwt.NewAuth,
	Runtime: kubernetes.NewRuntime,
	Broker: func(opts ...broker.Option) broker.Broker {
		return nats.NewBroker(broker.Addrs("nats-cluster"))
	},
	Registry: func(opts ...registry.Option) registry.Registry {
		return etcd.NewRegistry(registry.Addrs("etcd-cluster"))
	},
	Store:      cockroach.NewStore,
	AfterSetup: setupJWTRules,
}

// Service is the default for any services run
var Service = &Profile{}

// Test profile is used for the go test suite
var Test = &Profile{
	Auth:     noop.NewAuth,
	Store:    mem.NewStore,
	Registry: memReg.NewRegistry,
}
