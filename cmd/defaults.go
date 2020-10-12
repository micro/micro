package cmd

import (
	mucpNet "github.com/micro/micro/v3/service/network/mucp"
	"github.com/micro/micro/v3/service/auth"
	authSrv "github.com/micro/micro/v3/service/auth/client"
	"github.com/micro/micro/v3/service/network"
)

// setupDefaults sets the default auth, broker etc implementations incase they werent configured by
// a profile. The default implementations are always the RPC implementations.
func setupDefaults() {
	if auth.DefaultAuth == nil {
		auth.DefaultAuth = authSrv.NewAuth()
	}
	if network.DefaultNetwork == nil {
		network.DefaultNetwork = mucpNet.NewNetwork()
	}
}
