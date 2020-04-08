// Package cli implements the `micro store` subcommands
// for example:
//   micro store snapshot
//   micro store restore
package cli

import (
	"github.com/micro/cli/v2"
)

// CommonFlags are flags common to cli commands snapshot and restore
var CommonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "backend",
		Usage:   "Set the backend for the micro store",
		EnvVars: []string{"MICRO_STORE_BACKEND"},
		Value:   "memory",
	},
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
