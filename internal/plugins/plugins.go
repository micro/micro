// Package plugins includes the plugins we want to load
package plugins

import (
	"github.com/micro/go-micro/v2/config/cmd"

	// import specific plugins
	k8sRuntime "github.com/micro/go-micro/v2/runtime/kubernetes"
	cfStore "github.com/micro/go-micro/v2/store/cloudflare"
	ckStore "github.com/micro/go-micro/v2/store/cockroach"
	memStore "github.com/micro/go-micro/v2/store/file"
	fileStore "github.com/micro/go-micro/v2/store/memory"
)

func init() {
	// TODO: make it so we only have to import them
	cmd.DefaultRuntimes["kubernetes"] = k8sRuntime.NewRuntime
	cmd.DefaultStores["cloudflare"] = cfStore.NewStore
	cmd.DefaultStores["cockroach"] = ckStore.NewStore
	cmd.DefaultStores["file"] = fileStore.NewStore
	cmd.DefaultStores["memory"] = memStore.NewStore
}
