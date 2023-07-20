// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
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
	"micro.dev/v4/service/network"
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
	"micro.dev/v4/util/wrapper"

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

var (
	once sync.Once
)

// profiles which when called will configure micro to run in that environment
var profiles = map[string]*Profile{
	// built in profiles
	"client": Client,
	"server": Server,
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

		// set the network
		client.DefaultClient.Init(
			client.Network(network.Address),
		)

		return nil
	},
}

var Server = &Profile{
	Name: "server",
	Setup: func(ctx *cli.Context) error {
		// catch all
		SetupDefaults()

		// set auth
		auth.DefaultAuth = jwt.NewAuth(auth.Issuer(ctx.String("namespace")))

		SetupRules()
		// setup jwt
		SetupJWT()

		if ctx.Args().Get(1) == "registry" {
			SetupRegistry(memory.NewRegistry())
		} else {
			// set the registry address
			registry.DefaultRegistry.Init(
				registry.Addrs("localhost:8000"),
			)

			SetupRegistry(registry.DefaultRegistry)
		}

		if ctx.Args().Get(1) == "broker" {
			SetupBroker(memBroker.NewBroker())
		} else {
			broker.DefaultBroker.Init(
				broker.Addrs("localhost:8003"),
			)
			SetupBroker(broker.DefaultBroker)
		}

		// set store
		store.DefaultStore = file.NewStore(file.WithDir(filepath.Join(user.Dir, "server", "store")))

		// set config
		SetupConfigSecretKey()
		config.DefaultConfig, _ = storeConfig.NewConfig(store.DefaultStore, "")

		// setup events
		var err error
		events.DefaultStream, err = memStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}
		events.DefaultStore = evStore.NewStore(
			evStore.WithStore(store.DefaultStore),
		)

		// set the store in the model
		model.DefaultModel = sql.NewModel()

		// use the local runtime, note: the local runtime is designed to run source code directly so
		// the runtime builder should NOT be set when using this implementation
		runtime.DefaultRuntime = local.NewRuntime()

		// set blob store
		store.DefaultBlobStore, err = file.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		// set user
		SetupAccount(ctx)

		return nil
	},
}

func SetupJWT() {
	// get public/private key
	privKey, pubKey, err := user.GetJWTCerts()
	if err != nil {
		logger.Fatalf("Error getting keys: %v", err)
	}

	auth.DefaultAuth.Init(
		auth.PublicKey(string(pubKey)),
		auth.PrivateKey(string(privKey)),
	)
}

func SetupAccount(ctx *cli.Context) {
	opts := auth.DefaultAuth.Options()

	// extract the account creds from options, these can be set by flags
	accID := opts.ID
	accSecret := opts.Secret
	issuer := ""

	if ctx != nil {
		issuer = ctx.String("namespace")
	}

	// if no credentials were provided, self generate an account
	if len(accID) == 0 || len(accSecret) == 0 {
		opts := []auth.GenerateOption{
			auth.WithType("service"),
			auth.WithScopes("service"),
			auth.WithIssuer(issuer),
		}

		acc, err := auth.Generate(uuid.New().String(), opts...)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Debugf("Auth [%v] Generated an auth account", auth.DefaultAuth.String())

		accID = acc.ID
		accSecret = acc.Secret
	}

	// generate the first token
	token, err := auth.Token(
		auth.WithCredentials(accID, accSecret),
		auth.WithExpiry(time.Hour),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Debugf("Generated %v for acc %s %s", token, accID, accSecret)

	// set the credentials and token in auth options
	auth.DefaultAuth.Init(
		auth.ClientToken(token),
		auth.Credentials(accID, accSecret),
	)
}

func SetupDefaults() {
	once.Do(func() {
		client.DefaultClient = grpcCli.NewClient()
		server.DefaultServer = grpcSvr.NewServer()

		// wrap the client
		client.DefaultClient = wrapper.AuthClient(client.DefaultClient)
		// wrap the server
		server.DefaultServer.Init(server.WrapHandler(wrapper.AuthHandler()))

		// setup broker/registry
		SetupBroker(brokerSrv.NewBroker())
		SetupRegistry(registrySrv.NewRegistry())

		// setup rpc implementations after the client is configured
		config.DefaultConfig = configSrv.NewConfig()
		auth.DefaultAuth = authSrv.NewAuth()
		events.DefaultStream = eventsSrv.NewStream()
		events.DefaultStore = eventsSrv.NewStore()
		store.DefaultStore = storeSrv.NewStore()
		store.DefaultBlobStore = storeSrv.NewBlobStore()
		runtime.DefaultRuntime = runtimeSrv.NewRuntime()
		model.DefaultModel = sql.NewModel()

		// use the internal network lookup
		client.DefaultClient.Init(
			client.Lookup(network.Lookup),
		)

		// set the registry and broker in the client and server
		client.DefaultClient.Init(
			client.Broker(broker.DefaultBroker),
			client.Registry(registry.DefaultRegistry),
		)
		server.DefaultServer.Init(
			server.Broker(broker.DefaultBroker),
			server.Registry(registry.DefaultRegistry),
		)
	})
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
