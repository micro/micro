// Package platform/cli is for platform specific commands that are not yet dynamically generated
package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/auth"
	cl "github.com/micro/go-micro/v3/client"
	clinamespace "github.com/micro/micro/v3/client/cli/namespace"
	clitoken "github.com/micro/micro/v3/client/cli/token"
	cliutil "github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/report"
	pb "github.com/micro/micro/v3/platform/proto/signup"
	muclient "github.com/micro/micro/v3/service/client"
	"golang.org/x/crypto/ssh/terminal"
)

// Signup flow for the Micro Platform
func Signup(ctx *cli.Context) error {
	email := ctx.String("email")
	env := cliutil.GetEnv(ctx)
	reader := bufio.NewReader(os.Stdin)

	// no email specified
	if len(email) == 0 {
		// get email from prompt
		fmt.Print("Enter email address: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}

	// send a verification email to the user
	signupService := pb.NewSignupService("go.micro.service.signup", muclient.DefaultClient)
	_, err := signupService.SendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{
		Email: email,
	}, cl.WithRequestTimeout(10*time.Second))
	if err != nil {
		fmt.Printf("Error sending email during signup: %s\n", err)
		report.Errorf(ctx, "%v: Error sending email during signup: %s", email, err)
		os.Exit(1)
	}

	fmt.Print("We have sent you an email with a one time password. Please enter here: ")
	otp, _ := reader.ReadString('\n')
	otp = strings.TrimSpace(otp)

	// verify the email and password entered
	rsp, err := signupService.Verify(context.TODO(), &pb.VerifyRequest{
		Email: email,
		Token: otp,
	}, cl.WithRequestTimeout(10*time.Second))
	if err != nil {
		fmt.Printf("Error verifying: %s\n", err)
		report.Errorf(ctx, "%v: Error verifying: %s", email, err)
		os.Exit(1)
	}

	// Already registered users can just get logged in.
	tok := rsp.AuthToken
	if rsp.AuthToken != nil {

		err = clinamespace.Add(rsp.Namespace, env.Name)
		if err != nil {
			return err
		}
		err = clinamespace.Set(rsp.Namespace, env.Name)
		if err != nil {
			return err
		}
		if err := clitoken.Save(env.Name, &auth.Token{
			AccessToken:  tok.AccessToken,
			RefreshToken: tok.RefreshToken,
			Expiry:       time.Unix(tok.Expiry, 0),
		}); err != nil {
			return err
		}
		fmt.Println("Successfully logged in.")
		report.Success(ctx, email)
		return nil
	}

	// For users who don't have an account yet, this flow will proceed

	password := ctx.String("password")
	if len(password) == 0 {
		for {
			fmt.Print("Please enter your password: ")
			bytePw, _ := terminal.ReadPassword(int(syscall.Stdin))
			pw := string(bytePw)
			pw = strings.TrimSpace(pw)
			fmt.Println()

			fmt.Print("Please verify your password: ")
			bytePwVer, _ := terminal.ReadPassword(int(syscall.Stdin))
			pwVer := string(bytePwVer)
			pwVer = strings.TrimSpace(pwVer)
			fmt.Println()

			if pw != pwVer {
				fmt.Println("Passwords do not match. Please try again.")
				continue
			}
			password = pw
			break
		}
	}

	fmt.Printf("Please go to https://m3o.com/subscribe?email=%s and paste the acquired payment method id here: ", email)
	paymentMethodID, _ := reader.ReadString('\n')
	paymentMethodID = strings.TrimSpace(paymentMethodID)

	// complete the signup flow
	signupRsp, err := signupService.CompleteSignup(context.TODO(), &pb.CompleteSignupRequest{
		Email:           email,
		Token:           otp,
		PaymentMethodID: paymentMethodID,
		Secret:          password,
	}, cl.WithRequestTimeout(30*time.Second))
	if err != nil {
		fmt.Printf("Error completing signup: %s\n", err)
		report.Errorf(ctx, "Error completing signup: %s", err)
		os.Exit(1)
	}

	tok = signupRsp.AuthToken
	if err := clinamespace.Add(signupRsp.Namespace, env.Name); err != nil {
		fmt.Printf("Error adding namespace: %s\n", err)
		report.Errorf(ctx, "Error adding namespace: %s", err)
		os.Exit(1)
	}

	if err := clinamespace.Set(signupRsp.Namespace, env.Name); err != nil {
		fmt.Printf("Error setting namespace: %s\n", err)
		report.Errorf(ctx, "Error setting namespace: %s", err)
		os.Exit(1)
	}

	if err := clitoken.Save(env.Name, &auth.Token{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       time.Unix(tok.Expiry, 0),
	}); err != nil {
		fmt.Printf("Error saving token: %s\n", err)
		report.Errorf(ctx, "Error saving token: %s", err)
		os.Exit(1)
	}

	// the user has now signed up and logged in
	// @todo save the namespace from the last call and use that.
	fmt.Println("Successfully logged in.")
	report.Success(ctx, email)
	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:        "signup",
		Usage:       "Signup to the Micro Platform",
		Description: "Enables signup to the Micro Platform which can then be accessed via `micro env set platform` and `micro login`",
		Action:      Signup,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "email",
				Usage: "Email address to use for signup",
			},
			// In fact this is only here currently to help testing
			// as the signup flow can't be automated yet.
			// The testing breaks because we take the password
			// with the `terminal` package that makes input invisible.
			// That breaks tests though so password flag is used to get around tests.
			// @todo maybe payment method token and email sent verification
			// code should also be invisible. Problem for an other day.
			&cli.StringFlag{
				Name:  "password",
				Usage: "Password to use for login. If not provided, will be asked for during login. Useful for automated scripts",
			},
		},
	})
}
