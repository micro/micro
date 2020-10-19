// Package ci is for continuous integration testing
package ci

import (
	"github.com/micro/micro/v3/service/runtime/local"
	"github.com/micro/micro/plugin/etcd/v3"
	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/auth/jwt"
	"github.com/micro/micro/v3/service/broker/http"
	"github.com/micro/micro/v3/service/config"
	storeConfig "github.com/micro/micro/v3/service/config/store"
	"github.com/micro/micro/v3/service/events"
	evStore "github.com/micro/micro/v3/service/events/store"
	memStream "github.com/micro/micro/v3/service/events/stream/memory"
	"github.com/micro/micro/v3/service/logger"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/file"
	"github.com/urfave/cli/v2"
)

func init() {
	profile.Register("ci", Profile)
}

// CI profile to use for CI tests
var Profile = &profile.Profile{
	Name: "ci",
	Setup: func(ctx *cli.Context) error {
		auth.DefaultAuth = jwt.NewAuth()
		microRuntime.DefaultRuntime = local.NewRuntime()
		store.DefaultStore = file.NewStore()
		config.DefaultConfig, _ = storeConfig.NewConfig(store.DefaultStore, "")
		events.DefaultStream, _ = memStream.NewStream()
		events.DefaultStore = evStore.NewStore(evStore.WithStore(store.DefaultStore))
		profile.SetupBroker(http.NewBroker())
		profile.SetupRegistry(etcd.NewRegistry())
		profile.SetupJWT(ctx)
		profile.SetupConfigSecretKey(ctx)

		var err error
		store.DefaultBlobStore, err = file.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		return nil
	},
}
