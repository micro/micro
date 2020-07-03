// +build m3o

package test

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/token"
)

func TestM3oSignupFlow(t *testing.T) {
	trySuite(t, testM3oSignupFlow, retryCount)
}

func testM3oSignupFlow(t *t) {
	t.Parallel()

	serv := newServer(t)
	serv.launch()
	defer serv.close()

	envToConfigKey := map[string]string{
		"MICRO_STRIPE_API_KEY":           "micro.payments.stripe.api_key",
		"MICRO_SENDGRID_API_KEY":         "micro.signup.sendgrid.api_key",
		"MICRO_SENDGRID_TEMPLATE_ID":     "micro.signup.sendgrid.template_id",
		"MICRO_STRIPE_PLAN_ID":           "micro.signup.plan_id",
		"MICRO_EMAIL_FROM":               "micro.signup.email_from",
		"MICRO_TEST_ENV":                 "micro.signup.test_env",
		"MICRO_STRIPE_PAYMENT_METHOD_ID": "micro.signup.test_payment_method",
	}

	for envKey, configKey := range envToConfigKey {
		val := os.Getenv(envKey)
		if len(val) == 0 {
			t.Fatalf("'%v' flag is missing", envKey)
		}
		outp, err := exec.Command("micro", serv.envFlag(), "config", "set", configKey, val).CombinedOutput()
		if err != nil {
			t.Fatal(string(outp))
		}
	}

	outp, err := exec.Command("micro", serv.envFlag(), "run", "github.com/micro/services/signup").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	outp, err = exec.Command("micro", serv.envFlag(), "run", "github.com/micro/services/payments/provider/stripe").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	try("Find signup and stripe in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "stripe") || !strings.Contains(string(outp), "signup") {
			return outp, errors.New("Can't find sign or stripe in list")
		}
		return outp, err
	}, 50*time.Second)

	cmd := exec.Command("micro", serv.envFlag(), "login")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		outp, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(string(outp), err)
		}
		if !strings.Contains(string(outp), "Success") {
			t.Fatal(string(outp))
		}
		wg.Done()
	}()
	go func() {
		time.Sleep(15 * time.Second)
		cmd.Process.Kill()
	}()
	_, err = io.WriteString(stdin, "dobronszki@gmail.com\n")
	if err != nil {
		t.Fatal(err)
	}

	code := ""
	try("Find verification token in logs", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "10", "signup")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "Sending verification token") {
			return outp, errors.New("Output does not contain expected")
		}
		for _, line := range strings.Split(string(outp), "\n") {
			if strings.Contains(line, "Sending verification token") {
				code = strings.Split(line, "'")[1]
			}
		}
		return outp, nil
	}, 50*time.Second)

	_, err = io.WriteString(stdin, code+"\n")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)

	params := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String("4242424242424242"),
			ExpMonth: stripe.String("12"),
			ExpYear:  stripe.String("2021"),
			CVC:      stripe.String("123"),
		},
	}
	tok, err := token.New(params)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.WriteString(stdin, tok.ID+"\n")
	if err != nil {
		t.Fatal(err)
	}
	// Don't wait if a test is already failed, this is a quirk of the
	// test framework @todo fix this quirk
	if t.failed {
		return
	}
	wg.Wait()
}
