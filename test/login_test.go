// +build kind

package test

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestCorruptedTokenLogin checks that if we corrupt the token we successfully reset the config and clear the token
// to allow the user to login again rather than leave them in a state of limbo where they have to munge the config
// themselves
func TestCorruptedTokenLogin(t *testing.T) {
	// @todo this test is disabled now because it was
	// built on internal assumptions that no longer hold and not so easy to access anymore
	// TrySuite(t, testCorruptedLogin, retryCount)
}

func testCorruptedLogin(t *T) {
	serv := NewServer(t)
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	t.Parallel()

	// get server command
	cmd := serv.Command()

	outp, _ := cmd.Exec("status")
	unauthorizedMessage := "account is required"
	if !strings.Contains(string(outp), unauthorizedMessage) {
		t.Fatalf("Call should need authorization")
	}
	outp, _ = cmd.Exec("login", "--email", "admin", "--password", "micro")
	if !strings.Contains(string(outp), "Successfully logged in.") {
		t.Fatalf("Login failed: %s", outp)
	}
	outp, _ = cmd.Exec("status")
	if string(outp) != "" {
		t.Fatalf("Call should receive no output: %s", outp)
	}
	// munge token
	tok, err := cmd.Exec("user", "config", "get", "micro.auth."+serv.Env()+".refresh-token")
	if err != nil {
		t.Fatalf("Error getting refresh token value %s", err)
	}
	if _, err := cmd.Exec("user", "config", "set", "micro.auth."+serv.Env()+".refresh-token", strings.TrimSpace(string(tok))+"a"); err != nil {
		t.Fatalf("Error setting refresh token value %s", err)
	}
	if _, err := cmd.Exec("user", "config", "set", "micro.auth."+serv.Env()+".expiry", fmt.Sprintf("%d", time.Now().Add(-1*time.Hour).Unix())); err != nil {
		t.Fatalf("Error getting refresh token expiry %s", err)
	}

	outp, _ = cmd.Exec("status")
	if !strings.Contains(string(outp), unauthorizedMessage) {
		t.Fatalf("Call should have failed: %s", outp)
	}
	outp, _ = cmd.Exec("login", "--email", "admin", "--password", "micro")
	if !strings.Contains(string(outp), "Successfully logged in.") {
		t.Fatalf("Login failed: %s", outp)
	}

}
