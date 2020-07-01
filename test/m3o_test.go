// +build m3o

package test

import (
	"os"
	"os/exec"
	"testing"
)

func TestM3oSignupFlow(t *testing.T) {
	trySuite(t, testM3oSignupFlow, retryCount)
}

func testM3oSignupFlow(t *t) {
	t.Parallel()
	p, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(p) == 0 {
		t.Fatal("Git is not available")
		return
	}
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	envToConfigKey := map[string]string{
		"MICRO_STRIPE_API_KEY":       "micro.payments.stripe.api_key",
		"MICRO_SENDGRID_API_KEY":     "micro.signup.sendgrid.api_key",
		"MICRO_SENDGRID_TEMPLATE_ID": "micro.signup.sendgrid.template_id",
		"MICRO_STRIPE_PLAN_ID":       "micro.signup.sendgrid.template_id",
		"MICRO_STRIPE_PLAN_ID":       "micro.signup.plan_id",
		"MICRO_EMAIL_FROM":           "micro.signup.email_from",
		"MICRO_TEST_ENV", "micro.test_env",
	}

	for envKey, configKey := range envToConfigKey {
		val := os.Getenv(envKey)
		if len(val) == 0 {
			t.Fatalf("'%v' flag is missing", envKey)
		}
		outp, err = exec.Command("micro", serv.envFlag(), "config", "set", configKey, val.CombinedOutput()
		if err != nil {
			t.Fatal(string(outp))
		}
	}
}
