// Package cli implements the `micro store` subcommands
// for example:
//
//	micro store snapshot
//	micro store restore
//	micro store sync
package cli

import (
	"github.com/micro/micro/v5/cmd"
	"github.com/micro/micro/v5/util/helper"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd.Register(&cli.Command{
		Name:   "store",
		Usage:  "Commands for accessing the store",
		Action: helper.UnexpectedSubcommand,
		Subcommands: []*cli.Command{
			{
				Name:      "read",
				Usage:     "read a record from the store",
				UsageText: `micro store read [options] key`,
				Action:    read,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "database",
						Aliases: []string{"d"},
						Usage:   "database to write to",
						Value:   "micro",
					},
					&cli.StringFlag{
						Name:    "table",
						Aliases: []string{"t"},
						Usage:   "table to write to",
						Value:   "micro",
					},
					&cli.BoolFlag{
						Name:    "prefix",
						Aliases: []string{"p"},
						Usage:   "read prefix",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "suffix",
						Aliases: []string{"s"},
						Usage:   "read suffix",
						Value:   false,
					},
					&cli.UintFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "list limit",
					},
					&cli.StringFlag{
						Name:  "order",
						Usage: "Set the order of records e.g asc or desc",
					},
					&cli.UintFlag{
						Name:    "offset",
						Aliases: []string{"o"},
						Usage:   "list offset",
					},
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "show keys and headers (only values shown by default)",
						Value:   false,
					},
					&cli.StringFlag{
						Name:  "output",
						Usage: "output format (json, table)",
						Value: "table",
					},
				},
			},
			{
				Name:      "list",
				Usage:     "list all keys from a store",
				UsageText: `micro store list [options]`,
				Action:    list,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "database",
						Aliases: []string{"d"},
						Usage:   "database to list from",
						Value:   "micro",
					},
					&cli.StringFlag{
						Name:    "table",
						Aliases: []string{"t"},
						Usage:   "table to write to",
						Value:   "micro",
					},
					&cli.StringFlag{
						Name:  "output",
						Usage: "output format (json)",
					},
					&cli.StringFlag{
						Name:  "order",
						Usage: "Set the order of records e.g asc or desc",
					},
					&cli.BoolFlag{
						Name:    "prefix",
						Aliases: []string{"p"},
						Usage:   "list prefix",
						Value:   false,
					},
					&cli.UintFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "list limit",
					},
					&cli.UintFlag{
						Name:    "offset",
						Aliases: []string{"o"},
						Usage:   "list offset",
					},
				},
			},
			{
				Name:      "write",
				Usage:     "write a record to the store",
				UsageText: `micro store write [options] key value`,
				Action:    write,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "expiry",
						Aliases: []string{"e"},
						Usage:   "expiry in time.ParseDuration format",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "database",
						Aliases: []string{"d"},
						Usage:   "database to write to",
						Value:   "micro",
					},
					&cli.StringFlag{
						Name:    "table",
						Aliases: []string{"t"},
						Usage:   "table to write to",
						Value:   "micro",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a key from the store",
				UsageText: `micro store delete [options] key`,
				Action:    delete,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "database",
						Usage: "database to delete from",
						Value: "micro",
					},
					&cli.StringFlag{
						Name:  "table",
						Usage: "table to delete from",
						Value: "micro",
					},
				},
			},
			{
				Name:   "databases",
				Usage:  "List all databases known to the store service",
				Action: databases,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "store",
						Usage: "store service to call",
						Value: "store",
					},
				},
			},
			{
				Name:   "tables",
				Usage:  "List all tables in the specified database known to the store service",
				Action: tables,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "store",
						Usage: "store service to call",
						Value: "store",
					},
					&cli.StringFlag{
						Name:    "database",
						Aliases: []string{"d"},
						Usage:   "database to list tables of",
						Value:   "micro",
					},
				},
			},
			{
				Name:   "snapshot",
				Usage:  "Back up a store",
				Action: snapshot,
				Flags: append(CommonFlags,
					&cli.StringFlag{
						Name:    "destination",
						Usage:   "Backup destination",
						Value:   "file:///tmp/store-snapshot",
						EnvVars: []string{"MICRO_SNAPSHOT_DESTINATION"},
					},
				),
			},
			{
				Name:   "sync",
				Usage:  "Copy all records of one store into another store",
				Action: sync,
				Flags:  SyncFlags,
			},
			{
				Name:   "restore",
				Usage:  "restore a store snapshot",
				Action: restore,
				Flags: append(CommonFlags,
					&cli.StringFlag{
						Name:  "source",
						Usage: "Backup source",
						Value: "file:///tmp/store-snapshot",
					},
				),
			},
		},
	})
}

// CommonFlags are flags common to cli commands snapshot and restore
var CommonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "nodes",
		Usage:   "Comma separated list of Nodes to pass to the store backend",
		EnvVars: []string{"MICRO_STORE_NODES"},
	},
	&cli.StringFlag{
		Name:    "database",
		Usage:   "Database option to pass to the store backend",
		EnvVars: []string{"MICRO_STORE_DATABASE"},
	},
	&cli.StringFlag{
		Name:    "table",
		Usage:   "Table option to pass to the store backend",
		EnvVars: []string{"MICRO_STORE_TABLE"},
	},
}
