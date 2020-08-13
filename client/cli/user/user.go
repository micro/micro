// Package user handles the user cli command
package user

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/micro/cli/v2"
	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/config"
	"github.com/micro/micro/v3/service/auth"
	pb "github.com/micro/micro/v3/service/auth/proto"
	"github.com/micro/micro/v3/service/client"
	"golang.org/x/crypto/ssh/terminal"
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
				{
					Name:  "set",
					Usage: "Set various user based properties, eg. password",
					Subcommands: []*cli.Command{
						{
							Name:   "password",
							Usage:  "Set password",
							Action: changePassword,
							Flags: []cli.Flag{
								&cli.StringFlag{
									Name:  "old-password",
									Usage: "Existing password, the one that is used currently.",
								},
								&cli.StringFlag{
									Name:  "new-password",
									Usage: "New password you want to set.",
								},
							},
						},
					},
				},
			},
		},
	)
}

// get current user settings
func changePassword(ctx *cli.Context) error {
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

	oldPassword := ctx.String("old-password")
	newPassword := ctx.String("new-password")

	if len(oldPassword) == 0 {
		fmt.Print("Enter current password: ")
		bytePw, _ := terminal.ReadPassword(int(syscall.Stdin))
		pw := string(bytePw)
		pw = strings.TrimSpace(pw)
		fmt.Println()
		oldPassword = pw
	}

	if len(newPassword) == 0 {
		for {
			fmt.Print("Enter a new password: ")
			bytePw, _ := terminal.ReadPassword(int(syscall.Stdin))
			pw := string(bytePw)
			pw = strings.TrimSpace(pw)
			fmt.Println()

			fmt.Print("Verify your password: ")
			bytePwVer, _ := terminal.ReadPassword(int(syscall.Stdin))
			pwVer := string(bytePwVer)
			pwVer = strings.TrimSpace(pwVer)
			fmt.Println()

			if pw != pwVer {
				fmt.Println("Passwords do not match. Please try again.")
				continue
			}
			newPassword = pw
			break
		}
	}

	// Inspect the token
	acc, err := auth.Inspect(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	accountService := pb.NewAccountsService("auth", client.DefaultClient)
	accountService.ChangePassword(context.TODO(), &pb.ChangePasswordRequest{
		Id:        acc.ID,
		OldSecret: oldPassword,
		NewSecret: newPassword,
	}, goclient.WithAuthToken())
	return nil
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
	acc, err := auth.Inspect(token)
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
	acc, err := auth.Inspect(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(acc.ID)
	return nil
}
