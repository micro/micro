package cli

import (
	"github.com/micro/cli/v2"
	"github.com/pkg/errors"
)

// Sync is the entrypoint for micro store sync
func Sync(ctx *cli.Context) error {
	from, to, err := makeStores(ctx)
	if err != nil {
		return errors.Wrap(err, "Sync")
	}

	keys, err := from.List()
	if err != nil {
		return errors.Wrapf(err, "couldn't list from store %s", from.String())
	}
	for _, k := range keys {
		r, err := from.Read(k)
		if err != nil {
			return errors.Wrapf(err, "couldn't read %s from store %s", k, from.String())
		}
		if len(r) != 1 {
			return errors.Errorf("received multiple records reading %s from %s", k, from.String())
		}
		err = to.Write(r[0])
		if err != nil {
			return errors.Wrapf(err, "couldn't write %s to store %s", k, to.String())
		}
	}
	return nil
}

// SyncFlags are the flags for micro store sync
var SyncFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "from-backend",
		Usage:   "Backend to sync from",
		EnvVars: []string{"MICRO_STORE_FROM"},
	},
	&cli.StringFlag{
		Name:    "from-nodes",
		Usage:   "Nodes to sync from",
		EnvVars: []string{"MICRO_STORE_FROM_NODES"},
	},
	&cli.StringFlag{
		Name:    "from-database",
		Usage:   "Database to sync from",
		EnvVars: []string{"MICRO_STORE_FROM_DATABASE"},
	},
	&cli.StringFlag{
		Name:    "from-table",
		Usage:   "Table to sync from",
		EnvVars: []string{"MICRO_STORE_FROM_TABLE"},
	},
	&cli.StringFlag{
		Name:    "to-backend",
		Usage:   "Backend to sync to",
		EnvVars: []string{"MICRO_STORE_TO"},
	},
	&cli.StringFlag{
		Name:    "to-nodes",
		Usage:   "Nodes to sync to",
		EnvVars: []string{"MICRO_STORE_TO_NODES"},
	},
	&cli.StringFlag{
		Name:    "to-database",
		Usage:   "Database to sync to",
		EnvVars: []string{"MICRO_STORE_TO_DATABASE"},
	},
	&cli.StringFlag{
		Name:    "to-table",
		Usage:   "Table to sync to",
		EnvVars: []string{"MICRO_STORE_TO_TABLE"},
	},
}
