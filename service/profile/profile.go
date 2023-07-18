// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"micro.dev/v4/service/auth"
	"micro.dev/v4/service/auth/jwt"
	"micro.dev/v4/service/broker"
	memBroker "micro.dev/v4/service/broker/memory"
	"micro.dev/v4/service/client"
	"micro.dev/v4/service/config"
	storeConfig "micro.dev/v4/service/config/store"
	"micro.dev/v4/service/events"
	evStore "micro.dev/v4/service/events/store"
	memStream "micro.dev/v4/service/events/stream/memory"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/model"
	"micro.dev/v4/service/model/sql"
	"micro.dev/v4/service/registry"
	"micro.dev/v4/service/registry/memory"
	"micro.dev/v4/service/router"
	regRouter "micro.dev/v4/service/router/registry"
	"micro.dev/v4/service/runtime"
	"micro.dev/v4/service/runtime/local"
	"micro.dev/v4/service/server"
	"micro.dev/v4/service/store"
	"micro.dev/v4/service/store/file"
	inAuth "micro.dev/v4/util/auth"
	"micro.dev/v4/util/user"

	authSrv "micro.dev/v4/service/auth/client"
	brokerSrv "micro.dev/v4/service/broker/client"
	grpcCli "micro.dev/v4/service/client/grpc"
	configSrv "micro.dev/v4/service/config/client"
	eventsSrv "micro.dev/v4/service/events/client"
	registrySrv "micro.dev/v4/service/registry/client"
	runtimeSrv "micro.dev/v4/service/runtime/client"
	grpcSvr "micro.dev/v4/service/server/grpc"
	storeSrv "micro.dev/v4/service/store/client"
)

// profiles which when called will configure micro to run in that environment
var profiles = map[string]*Profile{
	// built in profiles
	"client":  Client,
	"service": Service,
	"server":  Server,
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
	Name: "client",
	Setup: func(ctx *cli.Context) error {
		SetupDefaults()
		return nil
	},
}

var Server = &Profile{
	Name: "server",
	Setup: func(ctx *cli.Context) error {
		// catch all
		SetupDefaults()

		// get public/private key
		privKey, pubKey, err := user.GetJWTCerts()
		if err != nil {
			logger.Fatalf("Error getting keys: %v", err)
		}

		// set auth
		auth.DefaultAuth = jwt.NewAuth(
			auth.PublicKey(string(pubKey)),
			auth.PrivateKey(string(privKey)),
		)

		// set broker
		SetupBroker(memBroker.NewBroker())

		// set store
		store.DefaultStore = file.NewStore(file.WithDir(filepath.Join(user.Dir, "server", "store")))

		// set config
		SetupConfigSecretKey()
		config.DefaultConfig, _ = storeConfig.NewConfig(store.DefaultStore, "")

		// setup events
		events.DefaultStream, err = memStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}
		events.DefaultStore = evStore.NewStore(
			evStore.WithStore(store.DefaultStore),
		)

		// set the store in the model
		model.DefaultModel = sql.NewModel()

		// set registry
		SetupRegistry(memory.NewRegistry())

		// use the local runtime, note: the local runtime is designed to run source code directly so
		// the runtime builder should NOT be set when using this implementation
		runtime.DefaultRuntime = local.NewRuntime()

		// set blob store
		store.DefaultBlobStore, err = file.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		// set jwt
		SetupRules()

		return nil
	},
}

// Service is the default for any services run
var Service = &Profile{
	Name: "service",
	Setup: func(ctx *cli.Context) error {
		SetupDefaults()
		return nil
	},
}

func SetupDefaults() {
	client.DefaultClient = grpcCli.NewClient()
	server.DefaultServer = grpcSvr.NewServer()

	// setup rpc implementations after the client is configured
	auth.DefaultAuth = authSrv.NewAuth()
	broker.DefaultBroker = brokerSrv.NewBroker()
	config.DefaultConfig = configSrv.NewConfig()
	events.DefaultStream = eventsSrv.NewStream()
	events.DefaultStore = eventsSrv.NewStore()
	registry.DefaultRegistry = registrySrv.NewRegistry()
	store.DefaultStore = storeSrv.NewStore()
	store.DefaultBlobStore = storeSrv.NewBlobStore()
	runtime.DefaultRuntime = runtimeSrv.NewRuntime()
	model.DefaultModel = sql.NewModel()
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

// SetupRules configures the default internal system rules
func SetupRules() {
	for _, rule := range inAuth.SystemRules {
		if err := auth.DefaultAuth.Grant(rule); err != nil {
			logger.Fatal("Error creating default rule: %v", err)
		}
	}
}

func SetupConfigSecretKey() {
	k, err := user.GetConfigSecretKey()
	if err != nil {
		logger.Fatal("Error getting config secret: %v", err)
	}
	os.Setenv("MICRO_CONFIG_SECRET_KEY", k)
}
