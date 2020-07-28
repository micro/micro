// Package user providers the micro user cli command
package cli

import (
	"fmt"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/config"
	muauth "github.com/micro/micro/v3/service/auth"
)

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "user",
			Usage:  "Print the current logged in user",
			Action: user,
			Subcommands: []*cli.Command{
				// config as a sub command,
				{
					Name:        "config",
					Usage:       "{set, get, delete} [key] [value]",
					Description: "Manage user related config like id, token, namespace, etc",
					Action:      current,
					Subcommands: config.Commands,
				},
				{
					Name:   "token",
					Usage:  "Get the current user token",
					Action: getToken,
				},
				{
					Name:   "namespace",
					Usage:  "Get the current namespace",
					Action: getNamespace,
				},
			},
		},
	)
}

// get current user settings
func current(ctx *cli.Context) error {
	env, err := config.Get("env")
	if err != nil || len(env) == 0 {
		env = "n/a"
	}

	ns, err := config.Get("namespaces", env, "current")
	if err != nil || len(ns) == 0 {
		ns = "n/a"
	}

	token, err := config.Get("micro", "auth", env, "token")
	if err != nil {
		return err
	}

	id := "n/a"

	// Inspect the token
	acc, err := muauth.DefaultAuth.Inspect(token)
	if err == nil {
		id = acc.ID
	}

	fmt.Println("user:", id)
	fmt.Println("namespace:", ns)
	fmt.Println("environment:", env)
	return nil
}

// get token for current env
func getToken(ctx *cli.Context) error {
	env, err := config.Get("env")
	if err != nil {
		return err
	}
	token, err := config.Get("micro", "auth", env, "token")
	if err != nil {
		return err
	}
	fmt.Println(token)
	return nil
}

// get namespace in current env
func getNamespace(ctx *cli.Context) error {
	env, err := config.Get("env")
	if err != nil {
		return err
	}
	namespace, err := config.Get("namespaces", env, "current")
	if err != nil {
		return err
	}
	fmt.Println(namespace)
	return nil
}

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
	acc, err := muauth.DefaultAuth.Inspect(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(acc.ID)
	return nil
}
