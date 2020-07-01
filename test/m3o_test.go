// +build m3o

package test

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
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
	}, 40*time.Second)
}
