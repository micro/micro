// +build integration

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/micro/micro/v2/client/cli/token"
)

func TestServerAuth(t *testing.T) {
	trySuite(t, testServerAuth, retryCount)
}

func testServerAuth(t *t) {
	t.Parallel()
	serv := newServer(t)
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	basicAuthSuite(serv, t)
}

func TestServerAuthJWT(t *testing.T) {
	trySuite(t, testServerAuthJWT, retryCount)
}

func testServerAuthJWT(t *t) {
	t.Parallel()
	serv := newServer(t, options{
		auth: "jwt",
	})
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	basicAuthSuite(serv, t)
}

func basicAuthSuite(serv testServer, t *t) {
	login(serv, t, "default", "password")

	// Execute first command in read to wait for store service
	// to start up
	if err := try("Calling micro auth list accounts", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "auth", "list", "accounts")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "admin") ||
			!strings.Contains(string(outp), "default") {
			return outp, fmt.Errorf("Output should contain default admin account")
		}
		return outp, nil
	}, 15*time.Second); err != nil {
		return
	}

	if err := try("Calling micro auth list rules", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "auth", "list", "rules")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "default") {
			return outp, fmt.Errorf("Output should contain default rule")
		}
		return outp, nil
	}, 8*time.Second); err != nil {
		return
	}

	if err := try("Try to get token with default account", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.auth", "Auth.Token", `{"id":"default","secret":"password"}`)
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]interface{}{}
		err = json.Unmarshal(outp, &rsp)
		token, ok := rsp["token"].(map[string]interface{})
		if !ok {
			return outp, errors.New("Can't find token")
		}
		if _, ok = token["access_token"].(string); !ok {
			return outp, fmt.Errorf("Can't find access token")
		}
		if _, ok = token["refresh_token"].(string); !ok {
			return outp, fmt.Errorf("Can't find access token")
		}
		if _, ok = token["refresh_token"].(string); !ok {
			return outp, fmt.Errorf("Can't find refresh token")
		}
		if _, ok = token["expiry"].(string); !ok {
			return outp, fmt.Errorf("Can't find access token")
		}
		return outp, nil
	}, 8*time.Second); err != nil {
		return
	}
}

func TestServerLoginJWT(t *testing.T) {
	trySuite(t, testServerAuthJWT, retryCount)
}

func testServerLoginJWT(t *t) {
	t.Parallel()
	serv := newServer(t, options{
		auth: "jwt",
	})
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	login(serv, t, "default", "password")
}

// Test bad tokens by messing up refresh token and trying to log in
// - this used to make even login fail which resulted in a UX deadlock
func TestServerBadTokenJWT(t *testing.T) {
	trySuite(t, testServerAuthJWT, retryCount)
}

func testServerBadTokenJWT(t *t) {
	t.Parallel()
	serv := newServer(t, options{
		auth: "jwt",
	})
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	login(serv, t, "default", "password")

	// Micro status should work
	if err := try("Micro status", t, func() ([]byte, error) {
		return exec.Command("micro", serv.envFlag(), "status").CombinedOutput()
	}, 3*time.Second); err != nil {
		return
	}

	// Modify rules so only logged in users can do anything

	// Add new rule that only lets logged in users do anything
	outp, err := exec.Command("micro", serv.envFlag(), "auth", "create", "rule", "--access=granted", "--scope='*'", "--resource='*:*:*'", "onlyloggedin").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	// Remove default rule
	outp, err = exec.Command("micro", serv.envFlag(), "auth", "delete", "rule", "default").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	// Micro status should still work as our user is already logged in
	if err := try("Micro status", t, func() ([]byte, error) {
		return exec.Command("micro", serv.envFlag(), "status").CombinedOutput()
	}, 3*time.Second); err != nil {
		return
	}

	// Now get the token and mess it up

	tok, err := token.Get(serv.envName())
	if err != nil {
		t.Fatal(err)
		return
	}

	tok.AccessToken = ""
	tok.Expiry = time.Time{}
	tok.RefreshToken = "some-random-junk"

	err = token.Save(serv.envName(), tok)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Micro status should fail
	if err := try("Micro status", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "status").CombinedOutput()
		if err == nil {
			return outp, errors.New("Micro status should fail")
		}
		return outp, err
	}, 3*time.Second); err != nil {
		return
	}

	login(serv, t, "default", "password")

	// Micro status should still work again after login
	if err := try("Micro status", t, func() ([]byte, error) {
		return exec.Command("micro", serv.envFlag(), "status").CombinedOutput()
	}, 3*time.Second); err != nil {
		return
	}
}

func TestServerLockdown(t *testing.T) {
	trySuite(t, testServerAuth, retryCount)
}

func testServerLockdown(t *t) {
	t.Parallel()
	serv := newServer(t)
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	lockdownSuite(serv, t)
}

func TestServerLockdownJWT(t *testing.T) {
	trySuite(t, testServerAuthJWT, retryCount)
}

func testServerLockdownJWT(t *t) {
	t.Parallel()
	serv := newServer(t, options{
		auth: "jwt",
	})
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	basicAuthSuite(serv, t)
}

func lockdownSuite(serv testServer, t *t) {
	// Execute first command in read to wait for store service
	// to start up
	if err := try("Calling micro auth list rules", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "auth", "list", "rules")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "default") {
			return outp, fmt.Errorf("Output should contain default rule")
		}
		return outp, nil
	}, 15*time.Second); err != nil {
		return
	}

	email := "me@email.com"
	pass := "mystrongpass"

	outp, err = exec.Command("micro", serv.envFlag(), "auth", "create", "account", "--secret", pass, "--scopes", "admin", email).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err := exec.Command("micro", serv.envFlag(), "auth", "create", "rule", "--access=granted", "--scope='*'", "--resource='*:*:*'", "onlyloggedin").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = exec.Command("micro", serv.envFlag(), "auth", "create", "rule", "--access=granted", "--scope=''", "authpublic").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = exec.Command("micro", serv.envFlag(), "auth", "delete", "rule", "default").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = exec.Command("micro", serv.envFlag(), "auth", "delete", "account", "default").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	if err := try("Listing rules should fail before login", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "auth", "list", "rules").CombinedOutput()
		if err == nil {
			return outp, errors.New("List rules should fail")
		}
		return outp, err
	}, 31*time.Second); err != nil {
		return
	}

	login(serv, t, "me@email.com", "mystrongpass")

	if err := try("Listing rules should pass after login", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "auth", "list", "rules").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "onlyloggedin") || !strings.Contains(string(outp), "authpublic") {
			return outp, errors.New("Can't find rules")
		}
		return outp, err
	}, 31*time.Second); err != nil {
		return
	}
}
