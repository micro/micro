package cli

import (
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/internal/helper"
	// imported specifically for signup
)

var (
	// ruleFlags are provided to commands which create or delete rules
	ruleFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "scope",
			Usage: "The scope to amend, e.g. 'user' or '*', leave blank to make public",
		},
		&cli.StringFlag{
			Name:  "resource",
			Usage: "The resource to amend in the format type:name:endpoint, e.g. service:go.micro.auth:*",
		},
		&cli.StringFlag{
			Name:  "access",
			Usage: "The access level, must be granted or denied",
			Value: "granted",
		},
		&cli.IntFlag{
			Name:  "priority",
			Usage: "The priority level, default is 0, the greater the number the higher the priority",
			Value: 0,
		},
	}
	// accountFlags are provided to the create account command
	accountFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "secret",
			Usage: "The account secret (password)",
		},
		&cli.StringSliceFlag{
			Name:  "scopes",
			Usage: "Comma seperated list of scopes to give the account",
		},
	}
)

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "auth",
			Action: helper.UnexpectedSubcommand,
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List auth resources",
					Subcommands: []*cli.Command{
						{
							Name:   "rules",
							Usage:  "List auth rules",
							Action: listRules,
						},
						{
							Name:   "accounts",
							Usage:  "List auth accounts",
							Action: listAccounts,
						},
					},
				},
				{
					Name:  "create",
					Usage: "Create an auth resource",
					Subcommands: []*cli.Command{
						{
							Name:   "rule",
							Usage:  "Create an auth rule",
							Flags:  ruleFlags,
							Action: createRule,
						},
						{
							Name:   "account",
							Usage:  "Create an auth account",
							Flags:  accountFlags,
							Action: createAccount,
						},
					},
				},
				{
					Name:  "delete",
					Usage: "Delete a auth resource",
					Subcommands: []*cli.Command{
						{
							Name:   "rule",
							Usage:  "Delete an auth rule",
							Flags:  ruleFlags,
							Action: deleteRule,
						},
						{
							Name:   "account",
							Usage:  "Delete an auth account",
							Flags:  ruleFlags,
							Action: deleteAccount,
						},
					},
				},
			},
		},
		&cli.Command{
			Name:        "login",
			Usage:       `Interactive login flow.`,
			Description: "Run 'micro login' for micro servers or 'micro login --otp' for the Micro Platform.",
			Action:      login,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "otp",
					Usage: "Login/signup with a One Time Password.",
				},
				&cli.StringFlag{
					Name:  "password",
					Usage: "Password to use for login. If not provided, will be asked for during login. Useful for automated scripts",
				},
				&cli.StringFlag{
					Name:  "email",
					Usage: "Email address to use for login",
				},
			},
		},
	)
}
