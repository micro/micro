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
)

func TestServerAuth(t *testing.T) {
	trySuite(t, testServerAuth, retryCount)
}

func testServerAuth(t *t) {
	t.Parallel()
	serv := newServer(t, withLogin())
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

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

	outp, err := exec.Command("micro", serv.envFlag(), "auth", "create", "account", "--secret", pass, "--scopes", "admin", email).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = exec.Command("micro", serv.envFlag(), "auth", "create", "rule", "--access=granted", "--scope='*'", "--resource='*:*:*'", "onlyloggedin").CombinedOutput()
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
