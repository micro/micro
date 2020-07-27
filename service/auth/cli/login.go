package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/token"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/report"
	platform "github.com/micro/micro/v2/platform/cli"
	muauth "github.com/micro/micro/v2/service/auth"
	"golang.org/x/crypto/ssh/terminal"
)

// login flow.
// For documentation of the flow please refer to https://github.com/micro/development/pull/223
func login(ctx *cli.Context) error {
	// assuming --otp go to platform.Signup
	if isOTP := ctx.Bool("otp"); isOTP {
		return platform.Signup(ctx)
	}

	// otherwise assume username/password login

	// get the environment
	env := util.GetEnv(ctx)
	// get the email address
	email := ctx.String("email")

	// email is blank
	if len(email) == 0 {
		fmt.Print("Enter email address: ")
		// read out the email from prompt if blank
		reader := bufio.NewReader(os.Stdin)
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}

	// clear tokens and try again
	if err := token.Remove(env.Name); err != nil {
		fmt.Printf("Error: %s\n", err)
		report.Errorf(ctx, "%v: Token remove: %v", email, err.Error())
		os.Exit(1)
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	password := ctx.String("password")
	if len(password) == 0 {
		pw, err := getPassword()
		if err != nil {
			return err
		}
		password = strings.TrimSpace(pw)
		fmt.Println()
	}
	tok, err := muauth.DefaultAuth.Token(auth.WithCredentials(email, password), auth.WithTokenIssuer(ns))
	if err != nil {
		fmt.Println(err)
		report.Errorf(ctx, "%v: Getting token: %v", email, err.Error())
		os.Exit(1)
	}
	token.Save(env.Name, tok)

	fmt.Println("Successfully logged in.")
	report.Success(ctx, email)
	return nil
}

// taken from https://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go
func getPassword() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}
