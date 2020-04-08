package snapshot

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/store"
	"github.com/pkg/errors"
)

func makeStore(ctx *cli.Context) (store.Store, error) {
	builtinStore, exists := cmd.DefaultStores[ctx.String("backend")]
	if !exists {
		return nil, errors.Errorf("store %s is not an implemented store - check your plugins", ctx.String("backend"))
	}
	s := builtinStore(
		store.Database(ctx.String("database")),
		store.Nodes(strings.Split(ctx.String("nodes"), ",")...),
		store.Database(ctx.String("database")),
		store.Table(ctx.String("table")),
	)
	if err := s.Init(); err != nil {
		return nil, errors.Wrapf(err, "Couldn't init % store", ctx.String("backend"))
	}
	return s, nil
}
