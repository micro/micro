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

	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/stripe/stripe-go/v71"
	stripe_client "github.com/stripe/stripe-go/v71/client"
)

func TestM3oSignupFlow(t *testing.T) {
	trySuite(t, testM3oSignupFlow, retryCount)
}

func testM3oSignupFlow(t *t) {
	t.Parallel()

	serv := newServer(t)
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	envToConfigKey := map[string]string{
		"MICRO_STRIPE_API_KEY":       "micro.payments.stripe.api_key",
		"MICRO_SENDGRID_API_KEY":     "micro.signup.sendgrid.api_key",
		"MICRO_SENDGRID_TEMPLATE_ID": "micro.signup.sendgrid.template_id",
		"MICRO_STRIPE_PLAN_ID":       "micro.signup.plan_id",
		"MICRO_EMAIL_FROM":           "micro.signup.email_from",
		"MICRO_TEST_ENV":             "micro.signup.test_env",
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

	outp, err := exec.Command("micro", serv.envFlag(), "run", getSrcString("M3O_INVITE_SVC", "github.com/micro/services/account/invite")).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	outp, err = exec.Command("micro", serv.envFlag(), "run", getSrcString("M3O_SIGNUP_SVC", "github.com/micro/services/signup")).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	outp, err = exec.Command("micro", serv.envFlag(), "run", getSrcString("M3O_STRIPE_SVC", "github.com/micro/services/payments/provider/stripe")).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	if err := try("Find signup and stripe in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "stripe") || !strings.Contains(string(outp), "signup") || !strings.Contains(string(outp), "invite") {
			return outp, errors.New("Can't find signup or stripe or invite in list")
		}
		return outp, err
	}, 70*time.Second); err != nil {
		return
	}

	time.Sleep(5 * time.Second)

	cmd := exec.Command("micro", serv.envFlag(), "login", "--otp")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		outp, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected an error for login but got none")
		} else if !strings.Contains(string(outp), "signup.notallowed") {
			t.Fatal(string(outp))
		}
		wg.Done()
	}()
	go func() {
		time.Sleep(20 * time.Second)
		cmd.Process.Kill()
	}()
	_, err = io.WriteString(stdin, "dobronszki@gmail.com\n")
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	if t.failed {
		return
	}

	outp, err = exec.Command("micro", serv.envFlag(), "call", "go.micro.service.account.invite", "Invite.Create", `{"email":"dobronszki@gmail.com"}`).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	cmd = exec.Command("micro", serv.envFlag(), "login", "--otp")
	stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		outp, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(string(outp), err)
		}
		if !strings.Contains(string(outp), "Success") {
			t.Fatal(string(outp))
		}
		ns, err := namespace.Get(serv.envName)
		if err != nil {
			t.Fatalf("Eror getting namespace: %v", err)
			return
		}
		defer func() {
			namespace.Remove(ns, serv.envName)
		}()
		if strings.Count(ns, "_") != 2 {
			t.Fatalf("Expected 2 underscores in namespace but namespace is: %v", ns)
			return
		}
		t.t.Logf("Namespace set is %v", ns)
	}()
	go func() {
		time.Sleep(20 * time.Second)
		cmd.Process.Kill()
	}()
	_, err = io.WriteString(stdin, "dobronszki@gmail.com\n")
	if err != nil {
		t.Fatal(err)
	}

	code := ""
	if err := try("Find verification token in logs", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "100", "signup")
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
	}, 50*time.Second); err != nil {
		return
	}

	t.Log("Code is ", code)
	if code == "" {
		t.Fatal("Code not found")
		return
	}
	_, err = io.WriteString(stdin, code+"\n")
	if err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(5 * time.Second)

	sc := stripe_client.New(os.Getenv("MICRO_STRIPE_API_KEY"), nil)
	pm, err := sc.PaymentMethods.New(
		&stripe.PaymentMethodParams{
			Card: &stripe.PaymentMethodCardParams{
				Number:   stripe.String("4242424242424242"),
				ExpMonth: stripe.String("7"),
				ExpYear:  stripe.String("2021"),
				CVC:      stripe.String("314"),
			},
			Type: stripe.String("card"),
		})
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = io.WriteString(stdin, pm.ID+"\n")
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

func getSrcString(envvar, dflt string) string {
	if env := os.Getenv(envvar); env != "" {
		return env
	}
	return dflt
}
