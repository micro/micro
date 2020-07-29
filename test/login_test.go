// +build kind

package test

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/micro/micro/v3/client/cli/token"
)

func TestCorruptedTokenLogin(t *testing.T) {
	trySuite(t, testCorruptedLogin, retryCount)
}

func testCorruptedLogin(t *t) {
	serv := newServer(t)
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	t.Parallel()

	outp, _ := exec.Command("micro", serv.envFlag(), "status").CombinedOutput()
	if !strings.Contains(string(outp), "Unauthorized") {
		t.Fatalf("Call should need authorization")
	}
	outp, _ = exec.Command("micro", serv.envFlag(), "login", "--email", serv.envName(), "--password", "password").CombinedOutput()
	if !strings.Contains(string(outp), "Successfully logged in.") {
		t.Fatalf("Login failed: %s", outp)
	}
	outp, _ = exec.Command("micro", serv.envFlag(), "status").CombinedOutput()
	if string(outp) != "" {
		t.Fatalf("Call should receive no output: %s", outp)
	}
	// munge token
	tok, _ := token.Get(serv.envName())
	tok.Expiry = time.Now().Add(-1 * time.Hour)
	tok.RefreshToken = tok.RefreshToken + "a"
	token.Save(serv.envName(), tok)

	outp, _ = exec.Command("micro", serv.envFlag(), "status").CombinedOutput()
	if !strings.Contains(string(outp), "Account can't be found for refresh token") {
		t.Fatalf("Call should have failed: %s", outp)
	}
	outp, _ = exec.Command("micro", serv.envFlag(), "login", "--email", serv.envName(), "--password", "password").CombinedOutput()
	if !strings.Contains(string(outp), "Successfully logged in.") {
		t.Fatalf("Login failed: %s", outp)
	}

}
