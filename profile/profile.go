// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"github.com/micro/go-micro/v3/auth/jwt"
	"github.com/micro/go-micro/v3/auth/noop"
	"github.com/micro/go-micro/v3/broker/http"
	"github.com/micro/go-micro/v3/broker/nats"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/registry/etcd"
	"github.com/micro/go-micro/v3/registry/mdns"
	"github.com/micro/go-micro/v3/registry/memory"
	"github.com/micro/go-micro/v3/router"
	"github.com/micro/go-micro/v3/runtime/kubernetes"
	"github.com/micro/go-micro/v3/runtime/local"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/store/cockroach"
	"github.com/micro/go-micro/v3/store/file"
	mem "github.com/micro/go-micro/v3/store/memory"

	inauth "github.com/micro/micro/v2/internal/auth"
	muauth "github.com/micro/micro/v2/service/auth"
	mubroker "github.com/micro/micro/v2/service/broker"
	muclient "github.com/micro/micro/v2/service/client"
	muregistry "github.com/micro/micro/v2/service/registry"
	murouter "github.com/micro/micro/v2/service/router"
	muruntime "github.com/micro/micro/v2/service/runtime"
	muserver "github.com/micro/micro/v2/service/server"
	mustore "github.com/micro/micro/v2/service/store"
)

// Profiles which when called will configure micro
var Profiles = map[string]Profile{
	"ci":       CI,
	"test":     Test,
	"local":    Local,
	"platform": Platform,
}

// Profile configures an environment
type Profile func()

// Test profile is used for the go test suite
var Test Profile = func() {
	muauth.DefaultAuth = noop.NewAuth()
	mustore.DefaultStore = mem.NewStore()
	setRegistry(memory.NewRegistry())
}

// CI profile to use for CI tests
var CI Profile = func() {
	muauth.DefaultAuth = jwt.NewAuth()
	mubroker.DefaultBroker = http.NewBroker()
	muruntime.DefaultRuntime = local.NewRuntime()
	mustore.DefaultStore = file.NewStore()
	setRegistry(etcd.NewRegistry())
	setupJWTRules()
}

// Local profile to use for the server locally
var Local Profile = func() {
	muauth.DefaultAuth = noop.NewAuth()
	mubroker.DefaultBroker = http.NewBroker()
	muruntime.DefaultRuntime = local.NewRuntime()
	mustore.DefaultStore = file.NewStore()
	setRegistry(mdns.NewRegistry())
	setupJWTRules()
}

// Platform profile to use for the server running in a
// production environment
var Platform Profile = func() {
	muauth.DefaultAuth = jwt.NewAuth()
	mubroker.DefaultBroker = nats.NewBroker()
	muruntime.DefaultRuntime = kubernetes.NewRuntime()
	mustore.DefaultStore = cockroach.NewStore()
	setRegistry(etcd.NewRegistry())
	setupJWTRules()
}

func setRegistry(reg registry.Registry) {
	muregistry.DefaultRegistry = reg
	murouter.DefaultRouter.Init(router.Registry(reg))
	muserver.DefaultServer.Init(server.Registry(reg))
	muclient.DefaultClient.Init(client.Registry(reg))
}

func setupJWTRules() {
	for _, rule := range inauth.SystemRules {
		if err := muauth.DefaultAuth.Grant(rule); err != nil {
			logger.Fatal("Error creating default rule: %v", err)
		}
	}
}
