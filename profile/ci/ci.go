// Package ci is for continuous integration testing
package ci

import (
	"github.com/micro/go-micro/v3/auth/jwt"
	"github.com/micro/go-micro/v3/broker/http"
	"github.com/micro/go-micro/v3/config"
	evStore "github.com/micro/go-micro/v3/events/store"
	memStream "github.com/micro/go-micro/v3/events/stream/memory"
	"github.com/micro/go-micro/v3/runtime/local"
	"github.com/micro/go-micro/v3/store/file"
	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"

	microAuth "github.com/micro/micro/v3/service/auth"
	microConfig "github.com/micro/micro/v3/service/config"
	microEvents "github.com/micro/micro/v3/service/events"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	microStore "github.com/micro/micro/v3/service/store"

	// external plugins
	"github.com/micro/go-plugins/registry/etcd/v3"
)

func init() {
	profile.Register("ci", Profile)
}

// CI profile to use for CI tests
var Profile = &profile.Profile{
	Name: "ci",
	Setup: func(ctx *cli.Context) error {
		microAuth.DefaultAuth = jwt.NewAuth()
		microRuntime.DefaultRuntime = local.NewRuntime()
		microStore.DefaultStore = file.NewStore()
		microConfig.DefaultConfig, _ = config.NewConfig()
		microEvents.DefaultStream, _ = memStream.NewStream()
		microEvents.DefaultStore = evStore.NewStore(evStore.WithStore(microStore.DefaultStore))
		profile.SetupBroker(http.NewBroker())
		profile.SetupRegistry(etcd.NewRegistry())
		profile.SetupJWT(ctx)

		var err error
		microStore.DefaultBlobStore, err = file.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		return nil
	},
}
