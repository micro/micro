package cmd

import (
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/client/grpc"
)

// setupDefaults sets the default auth, broker etc implementations incase they werent configured by
// a profile. The default implementations are always the RPC implementations.
func setupDefaults() {
	if client.DefaultClient == nil {
		client.DefaultClient = grpc.NewClient()
	}
}
