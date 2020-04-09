// Package plugins includes the plugins we want to load
package plugins

import (
	"github.com/micro/go-micro/v2/config/cmd"

	// import specific plugins
	cfStore "github.com/micro/go-micro/v2/store/cloudflare"
	ckStore "github.com/micro/go-micro/v2/store/cockroach"
	fileStore "github.com/micro/go-micro/v2/store/file"
	memStore "github.com/micro/go-micro/v2/store/memory"
)

func init() {
	// TODO: make it so we only have to import them
	cmd.DefaultStores["cloudflare"] = cfStore.NewStore
	cmd.DefaultStores["cockroach"] = ckStore.NewStore
	cmd.DefaultStores["file"] = fileStore.NewStore
	cmd.DefaultStores["memory"] = memStore.NewStore
}
