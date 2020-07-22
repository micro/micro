// Package user providers the micro user cli command
package cli

import (
	"fmt"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/internal/config"
)

// user returns info about the logged in user
func user(ctx *cli.Context) error {
	env := util.GetEnv(ctx)

	// Get the token from micro config
	token, err := config.Get("micro", "auth", env.Name, "token")
	if err != nil {
		fmt.Println("You are not logged in")
		os.Exit(1)
	}

	if len(token) == 0 {
		fmt.Println("You are not logged in")
		os.Exit(1)
	}

	// Inspect the token
	a, err := authFromContext(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	acc, err := a.Inspect(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(acc.ID)
	return nil
}

func init() {
	cmd.Register(
		&cli.Command{
			Name:        "user",
			Usage:       "Print the current logged in user",
			Action:      user,
			Subcommands: config.UserCommands,
		},
	)
}
