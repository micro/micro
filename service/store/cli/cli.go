// Package cli implements the `micro store` subcommands
// for example:
//   micro store snapshot
//   micro store restore
//   micro store sync
package cli

import (
	"github.com/micro/cli/v2"
)

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

// Commands for data storing
func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "read",
			Usage:     "read a record from the store",
			UsageText: `micro store read [options] key`,
			Action:    Read,
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
			Action:    List,
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
			Action:    Write,
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
			Action:    Delete,
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
			Action: Databases,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "store",
					Usage: "store service to call",
					Value: "go.micro.store",
				},
			},
		},
		{
			Name:   "tables",
			Usage:  "List all tables in the specified database known to the store service",
			Action: Tables,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "store",
					Usage: "store service to call",
					Value: "go.micro.store",
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
			Action: Snapshot,
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
			Action: Sync,
			Flags:  SyncFlags,
		},
		{
			Name:   "restore",
			Usage:  "restore a store snapshot",
			Action: Restore,
			Flags: append(CommonFlags,
				&cli.StringFlag{
					Name:  "source",
					Usage: "Backup source",
					Value: "file:///tmp/store-snapshot",
				},
			),
		},
	}
}
