// Package user handles the user cli command
package user

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/config"
	pb "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/urfave/cli/v2"
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
					Subcommands: []*cli.Command{
						{
							Name:   "set",
							Usage:  "Set namespace in the current environment",
							Action: setNamespace,
						},
					},
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
									Name:  "email",
									Usage: "Email to use for password change",
								},
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
	email := ctx.String("email")
	if len(email) == 0 {
		env := util.GetEnv(ctx)
		token, err := config.Get(config.Path("micro", "auth", env.Name, "token"))
		if err != nil {
			return err
		}

		// Inspect the token
		acc, err := auth.Inspect(token)
		if err != nil {
			fmt.Println("You are not logged in")
			return err
		}
		email = acc.ID
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
	ns, err := currNamespace(ctx)
	if err != nil {
		return err
	}

	accountService := pb.NewAccountsService("auth", client.DefaultClient)
	_, err = accountService.ChangeSecret(context.DefaultContext, &pb.ChangeSecretRequest{
		Id:        email,
		OldSecret: oldPassword,
		NewSecret: newPassword,
		Options:   &pb.Options{Namespace: ns},
	}, goclient.WithAuthToken())
	return err
}

// get current user settings
func current(ctx *cli.Context) error {
	env := util.GetEnv(ctx).Name
	if len(env) == 0 {
		env = "n/a"
	}

	ns, err := config.Get(config.Path("namespaces", env, "current"))
	if err != nil || len(ns) == 0 {
		ns = "n/a"
	}

	token, err := config.Get(config.Path("micro", "auth", env, "token"))
	if err != nil {
		return err
	}

	gitcreds, err := config.Get(config.Path("git", "credentials"))
	if err != nil {
		return err
	}
	if len(gitcreds) > 0 {
		gitcreds = "[hidden]"
	} else {
		gitcreds = "n/a"
	}

	id := "n/a"

	// Inspect the token
	acc, err := auth.Inspect(token)
	if err == nil {
		id = acc.Name
		if len(id) == 0 {
			id = acc.ID
		}
	}

	baseURL, _ := config.Get(config.Path("git", util.GetEnv(ctx).Name, "baseurl"))
	if len(baseURL) == 0 {
		baseURL, _ = config.Get(config.Path("git", "baseurl"))
	}
	if len(baseURL) == 0 {
		baseURL = "n/a"
	}

	fmt.Println("user:", id)
	fmt.Println("namespace:", ns)
	fmt.Println("environment:", env)
	fmt.Println("git.credentials:", gitcreds)
	fmt.Println("git.baseurl:", baseURL)
	return nil
}

// get token for current env
func getToken(ctx *cli.Context) error {
	env, err := config.Get("env")
	if err != nil {
		return err
	}
	token, err := config.Get(config.Path("micro", "auth", env, "token"))
	if err != nil {
		return err
	}
	fmt.Println(token)
	return nil
}

// get namespace in current env
func getNamespace(ctx *cli.Context) error {
	namespace, err := currNamespace(ctx)
	if err != nil {
		return err
	}
	fmt.Println(namespace)
	return nil
}

func currNamespace(ctx *cli.Context) (string, error) {
	env, err := config.Get("env")
	if err != nil {
		return "", err
	}
	namespace, err := config.Get(config.Path("namespaces", env, "current"))
	if err != nil {
		return "", err
	}
	return namespace, nil
}

// set namespace in current env
func setNamespace(ctx *cli.Context) error {
	if len(ctx.Args().First()) == 0 {
		return errors.New("No namespace specified")
	}
	env, err := config.Get("env")
	if err != nil {
		return err
	}
	return config.Set(config.Path("namespaces", env, "current"), ctx.Args().First())
}

// user returns info about the logged in user
func user(ctx *cli.Context) error {
	env := util.GetEnv(ctx)

	// Get the token from micro config
	token, err := config.Get(config.Path("micro", "auth", env.Name, "token"))
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
	// backward compatibility
	user := acc.Name
	if len(user) == 0 {
		user = acc.ID
	}
	fmt.Println(user)
	return nil
}
