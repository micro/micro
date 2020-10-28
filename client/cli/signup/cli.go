// Package signup is for the signup command backed by a signup service
package signup

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	clinamespace "github.com/micro/micro/v3/client/cli/namespace"
	clitoken "github.com/micro/micro/v3/client/cli/token"
	cliutil "github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/report"
	pb "github.com/micro/micro/v3/proto/signup"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

// Run execs flow for the signup
func Run(ctx *cli.Context) error {
	email := ctx.String("email")
	if ctx.Bool("recover") {
		signupService := pb.NewSignupService("signup", client.DefaultClient)
		_, err := signupService.Recover(context.DefaultContext, &pb.RecoverRequest{
			Email: email,
		}, client.WithRequestTimeout(10*time.Second))
		return err
	}

	env, err := cliutil.GetEnv(ctx)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)

	// no email specified
	if len(email) == 0 {
		// get email from prompt
		fmt.Print("Enter email address: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}

	// send a verification email to the user
	signupService := pb.NewSignupService("signup", client.DefaultClient)
	_, err = signupService.SendVerificationEmail(context.DefaultContext, &pb.SendVerificationEmailRequest{
		Email: email,
	}, client.WithRequestTimeout(10*time.Second))
	if err != nil {
		report.Errorf(ctx, "Error sending email to %v during signup: %s", email, err)
		return err
	}

	fmt.Print("Enter the OTP sent to your email address: ")
	otp, _ := reader.ReadString('\n')
	otp = strings.TrimSpace(otp)

	// verify the email and password entered
	rsp, err := signupService.Verify(context.DefaultContext, &pb.VerifyRequest{
		Email: email,
		Token: otp,
	}, client.WithRequestTimeout(10*time.Second))
	if err != nil {
		report.Errorf(ctx, "Error verifying %v: %s", email, err)
		return err
	}

	isJoining := false

	if ns := rsp.Namespaces; len(ns) > 0 {
		fmt.Printf("\nYou've been invited to the '%v' namespace.\nDo you want to join it or create your own? Please type \"own\" or \"join\": ", ns[0])

		for {
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			validAnswer := false
			switch answer {
			case "join":
				isJoining = true
				validAnswer = true
				break
			case "own":
				validAnswer = true
			default:
				fmt.Printf("Answer \"%v\" is invalid. Valid answers are: \"own\" or \"join\": ", answer)
			}
			if validAnswer {
				break
			}
		}
	}

	password := ctx.String("password")

	if len(password) == 0 {
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
			password = pw
			break
		}
	}

	// payment method id read from user input
	var paymentMethodID string

	// Only take payment method if not joining, ie. creating their own namespace
	// and M3O platform subscription
	if !isJoining {
		// print the message returned from the verification process
		if len(rsp.Message) > 0 {
			// print with space
			fmt.Printf("\n%s\n", rsp.Message)
		}

		// payment required
		if rsp.PaymentRequired {
			for {
				hasRsp, err := signupService.HasPaymentMethod(context.DefaultContext, &pb.HasPaymentMethodRequest{
					Token: otp,
				})
				if err == nil && hasRsp != nil && hasRsp.Has {
					break
				}
				time.Sleep(2 * time.Second)
			}
		}
	}

	// complete the signup flow
	signupNamespace := ""
	if isJoining && len(rsp.Namespaces) > 0 {
		signupNamespace = rsp.Namespaces[0]
	}
	signupRsp, err := signupService.CompleteSignup(context.DefaultContext, &pb.CompleteSignupRequest{
		Email:           email,
		Token:           otp,
		PaymentMethodID: paymentMethodID,
		Secret:          password,
		Namespace:       signupNamespace,
	}, client.WithRequestTimeout(30*time.Second))
	if err != nil {
		report.Errorf(ctx, "Error completing signup: %s", err)
		return err
	}

	tok := signupRsp.AuthToken
	if err := clinamespace.Add(signupRsp.Namespace, env.Name); err != nil {
		report.Errorf(ctx, "Error adding namespace: %s", err)
		return err
	}

	if err := clinamespace.Set(signupRsp.Namespace, env.Name); err != nil {
		report.Errorf(ctx, "Error setting namespace: %s", err)
		return err
	}

	if err := clitoken.Save(ctx, &auth.AccountToken{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       time.Unix(tok.Expiry, 0),
	}); err != nil {
		report.Errorf(ctx, "Error saving token: %s", err)
		return err
	}

	// the user has now signed up and logged in
	// @todo save the namespace from the last call and use that.
	fmt.Println("\nSignup complete! You're now logged in.")
	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:        "signup",
		Usage:       "Signup to the Micro Platform",
		Description: "Enables signup to the Micro Platform which can then be accessed via `micro env set platform` and `micro login`",
		Action:      Run,
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
			&cli.BoolFlag{
				Name:  "recover",
				Usage: "Emails you the namespaces you have access to. micro signup --recover --email=youremail@domain.com",
			},
		},
	})
}
