package token

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/internal/token"
)

func generate(ctx *cli.Context) error {
	email := ctx.String("email")
	pass := ctx.String("pass")

	if len(email) == 0 {
		// TODO return err
		fmt.Println("Email is blank (specify --email)")
		os.Exit(1)
	}

	// no pass first request it
	if len(pass) == 0 {
		// TODO return err
		if err := token.SendPass(email); err != nil {
			fmt.Println("Sending OTP pass failed:", err)
			os.Exit(1)
		}

		// wait for pass
		fmt.Print("Enter OTP: ")
		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanWords)
		s.Scan()
		pass = s.Text()
	}

	// generate
	t, err := token.Generate(email, pass)
	if err != nil {
		// TODO return err
		fmt.Println("Token generation failed:", err)
		os.Exit(1)
	}
	fmt.Println("Your token (set as MICRO_TOKEN_KEY env var or X-Micro-Token http header):")
	fmt.Println(t)
	return nil
}

func revoke(ctx *cli.Context) error {
	tk := ctx.String("token")
	if len(tk) == 0 {
		// TODO return err
		fmt.Println("Token is blank (specify --token)")
		os.Exit(1)
	}

	// revoke token
	if err := token.Revoke(tk); err != nil {
		// TODO return err
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Token revoked")
	return nil
}

func verify(ctx *cli.Context) error {
	tk := ctx.String("token")
	if len(tk) == 0 {
		// TODO return err
		fmt.Println("Token is blank (specify --token)")
		os.Exit(1)
	}

	// revoke token
	if err := token.Verify(tk); err != nil {
		// TODO return err
		fmt.Println("Verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("Token verified")
	return nil
}

func list(ctx *cli.Context) error {
	tokens, err := token.List()
	if err != nil {
		// TODO return err
		fmt.Println(err)
		os.Exit(1)
	}
	if len(tokens) == 0 {
		fmt.Println(`{}`)
		return nil
	}
	j, err := json.MarshalIndent(tokens, "", "\t")
	if err != nil {
		//TODO:	return err
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(j))
	return nil
}

func tokenCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "list",
			Usage:  "List tokens",
			Action: list,
		},
		{
			Name:   "generate",
			Usage:  "Generate an api token (specify --email)",
			Action: generate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Usage: "Email address to generate token for. OTP pass is sent to email",
				},
				&cli.StringFlag{
					Name:  "pass",
					Usage: "OTP pass sent in email",
				},
			},
		},
		{
			Name:   "revoke",
			Usage:  "Revoke an api token (specify --token)",
			Action: revoke,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "token",
					Usage: "Encoded token key to revoke",
				},
			},
		},
		{
			Name:   "verify",
			Usage:  "Verify an api token is valid (specify --token)",
			Action: verify,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "token",
					Usage: "Encoded token key to verify",
				},
			},
		},
	}
}

// Commands returns token commands
func Commands() []*cli.Command {
	return []*cli.Command{{
		Name:        "token",
		Usage:       "API token commands",
		Subcommands: tokenCommands(),
	}}
}
