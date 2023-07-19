package cli

import (
	"github.com/urfave/cli/v2"
	"micro.dev/v4/cmd"
	"micro.dev/v4/util/helper"
	// imported specifically for signup
)

var (
	// ruleFlags are provided to commands which create or delete rules
	ruleFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "scope",
			Usage: "the scope to amend, e.g. 'user' or '*', leave blank to make public",
		},
		&cli.StringFlag{
			Name:  "resource",
			Usage: "The resource to amend in the format type:name:endpoint, e.g. service:auth:*",
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
			Usage: "Comma separated list of scopes to give the account",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "The type of account e.g client, user, service",
		},
	}
	// accountFlags are provided to the create account command
	tokenFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "id",
			Usage: "The account id",
		},
		&cli.StringFlag{
			Name:  "secret",
			Usage: "The account secret (password)",
		},
		&cli.IntFlag{
			Name:  "expiry",
			Usage: "The expiry time in seconds. Default no expiry",
		},
	}
)

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "auth",
			Usage:  "Manage authentication, accounts and rules",
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
							Name:  "account",
							Usage: "Create an auth account",
							Flags: append(accountFlags, &cli.StringFlag{
								Name:  "namespace",
								Usage: "Namespace to use when creating the account",
							}),
							Action: createAccount,
						},
						{
							Name:  "token",
							Usage: "Create an auth token",
							Flags: append(tokenFlags, &cli.StringFlag{
								Name:  "namespace",
								Usage: "Namespace to use when creating the account",
							}),
							Action: createToken,
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
							Flags:  accountFlags,
							Action: deleteAccount,
						},
					},
				},
				{
					Name:  "update",
					Usage: "Update an auth resource",
					Subcommands: []*cli.Command{
						{
							Name:  "secret",
							Usage: "Update an auth account secret",
							Flags: append(accountFlags,
								&cli.StringFlag{
									Name:  "namespace",
									Usage: "Namespace to use when updating the account",
								},
								&cli.StringFlag{
									Name:  "old_secret",
									Usage: "The old account secret (password)",
								},
								&cli.StringFlag{
									Name:  "new_secret",
									Usage: "The new account secret (password)",
								},
							),
							Action: updateAccount,
						},
					},
				},
			},
		},
		&cli.Command{
			Name:        "login",
			Usage:       `Interactive login flow.`,
			Description: "Run 'micro login' for the server",
			Action:      login,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "password",
					Usage: "Password to use for login. If not provided, will be asked for during login. Useful for automated scripts",
				},
				&cli.StringFlag{
					Name:    "username",
					Usage:   "Username to use for login",
					Aliases: []string{"email"},
				},
			},
		},
		&cli.Command{
			Name:        "logout",
			Usage:       `Logout.`,
			Description: "Use 'micro logout' to delete your token in your current environment.",
			Action:      logout,
		},
	)
}
