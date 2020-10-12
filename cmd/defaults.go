package cmd

import (
	"github.com/micro/micro/v3/service/auth"
	authSrv "github.com/micro/micro/v3/service/auth/client"
)

// setupDefaults sets the default auth, broker etc implementations incase they werent configured by
// a profile. The default implementations are always the RPC implementations.
func setupDefaults() {
	if auth.DefaultAuth == nil {
		auth.DefaultAuth = authSrv.NewAuth()
	}
}
