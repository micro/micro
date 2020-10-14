package cli

import (
	"strings"

	"github.com/micro/micro/v3/service/store"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// makeStore is a helper function that creates a store for snapshot and restore
func makeStore(ctx *cli.Context) (store.Store, error) {
	builtinStore, err := getStore(ctx.String("store"))
	if err != nil {
		return nil, errors.Wrap(err, "makeStore")
	}
	s := builtinStore(
		store.Nodes(strings.Split(ctx.String("nodes"), ",")...),
		store.Database(ctx.String("database")),
		store.Table(ctx.String("table")),
	)
	if err := s.Init(); err != nil {
		return nil, errors.Wrapf(err, "Couldn't init %s store", ctx.String("store"))
	}
	return s, nil
}

// makeStores is a helper function that sets up 2 stores for sync
func makeStores(ctx *cli.Context) (store.Store, store.Store, error) {
	fromBuilder, err := getStore(ctx.String("from-backend"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "from store")
	}
	toBuilder, err := getStore(ctx.String("to-backend"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "to store")
	}
	from := fromBuilder(
		store.Nodes(strings.Split(ctx.String("from-nodes"), ",")...),
		store.Database(ctx.String("from-database")),
		store.Table(ctx.String("from-table")),
	)
	if err := from.Init(); err != nil {
		return nil, nil, errors.Wrapf(err, "from: couldn't init %s", ctx.String("from-backend"))
	}
	to := toBuilder(
		store.Nodes(strings.Split(ctx.String("to-nodes"), ",")...),
		store.Database(ctx.String("to-database")),
		store.Table(ctx.String("to-table")),
	)
	if err := to.Init(); err != nil {
		return nil, nil, errors.Wrapf(err, "to: couldn't init %s", ctx.String("to-backend"))
	}
	return from, to, nil
}

func getStore(s string) (func(...store.StoreOption) store.Store, error) {
	// builtinStore, exists := cmd.DefaultStores[s]
	// if !exists {
	// 	return nil, errors.Errorf("store %s is not an implemented store - check your plugins", s)
	// }
	// return builtinStore, nil
	return nil, nil
}
