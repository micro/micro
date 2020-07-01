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

	stripeAPIKey := os.Getenv("MICRO_STRIPE_API_KEY")
	if len(stripeAPIKey) == 0 {
		t.Fatal("Stripe api key is missing")
	}
	outp, err := exec.Command("micro", serv.envFlag(), "config", "set", "micro.payments.stripe.api_key", stripeAPIKey).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	sendgridAPIKey := os.Getenv("MICRO_SENDGRID_API_KEY")
	if len(sendgridAPIKey) == 0 {
		t.Fatal("Stripe api key is missing")
	}
	outp, err = exec.Command("micro", serv.envFlag(), "config", "set", "micro.signup.sendgrid.api_key", sendgridAPIKey).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	sendgridTemplateID := os.Getenv("MICRO_SENDGRID_TEMPLATE_ID")
	if len(sendgridTemplateID) == 0 {
		t.Fatal("Sendgrid template ID is missing")
	}
	outp, err = exec.Command("micro", serv.envFlag(), "config", "set", "micro.signup.sendgrid.template_id", sendgridTemplateID).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	stripePlanID := os.Getenv("MICRO_STRIPE_PLAN_ID")
	if len(stripePlanID) == 0 {
		t.Fatal("Stripe plan ID is missing")
	}
	outp, err = exec.Command("micro", serv.envFlag(), "config", "set", "micro.signup.plan_id", stripePlanID).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	emailFrom := os.Getenv("MICRO_EMAIL_FROM")
	if len(emailFrom) == 0 {
		t.Fatal("Email from address is missing")
	}
	outp, err = exec.Command("micro", serv.envFlag(), "config", "set", "micro.signup.email_from", emailFrom).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	testEnv := os.Getenv("MICRO_TEST_ENV")
	if len(testEnv) == 0 {
		t.Fatal("Test env flag is missing")
	}
	outp, err = exec.Command("micro", serv.envFlag(), "config", "set", "micro.test_env", testEnv).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}
}
