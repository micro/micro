// +build integration

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/internal/config"
)

// Test no default account generation in non-default namespaces
func TestNoDefaultAccount(t *testing.T) {
	TrySuite(t, testNoDefaultAccount, retryCount)
}

func testNoDefaultAccount(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	ns := "random-namespace"

	err := ChangeNamespace(cmd, serv.Env(), ns)
	if err != nil {
		t.Fatal(err)
		return
	}

	Try("Log in with user should fail", t, func() ([]byte, error) {
		out, err := serv.Command().Exec("login", "--email", "admin", "--password", "micro")
		if err == nil {
			return out, errors.New("Loggin in should error")
		}
		if strings.Contains(string(out), "Success") {
			return out, errors.New("Loggin in should error")
		}
		return out, nil
	}, 5*time.Second)

	Try("Run helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("run", "helloworld")
		if err == nil {
			return outp, errors.New("Run should error")
		}
		return outp, nil
	}, 5*time.Second)

	Try("Find helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err == nil {
			return outp, errors.New("Should not be able to do status")
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("Shouldn't find example helloworld in runtime")
		}
		return outp, nil
	}, 15*time.Second)
}

func TestPublicAPI(t *testing.T) {
	TrySuite(t, testPublicAPI, retryCount)
}

func testPublicAPI(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	outp, err := cmd.Exec("auth", "create", "account", "--secret", "micro", "--namespace", "random-namespace", "admin")
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	err = ChangeNamespace(cmd, serv.Env(), "random-namespace")
	if err != nil {
		t.Fatal(err)
		return
	}
	// login to admin account
	if err = Login(serv, t, "admin", "micro"); err != nil {
		t.Fatalf("Error logging in %s", err)
		return
	}

	if err := Try("Run helloworld", t, func() ([]byte, error) {
		return cmd.Exec("run", "./services/helloworld")
	}, 5*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if !statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("Can't find example helloworld in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Call helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("helloworld", "--name=joe")
		if err != nil {
			outp1, _ := cmd.Exec("logs", "helloworld")
			return append(outp, outp1...), err
		}
		if !strings.Contains(string(outp), "Msg") {
			return outp, err
		}
		return outp, err
	}, 90*time.Second); err != nil {
		return
	}

	if err := Try("curl helloworld", t, func() ([]byte, error) {
		bod, rsp, err := curl(serv, "random-namespace", "helloworld?name=Jane")
		if rsp == nil {
			return []byte(bod), fmt.Errorf("helloworld should have response, err: %v", err)
		}
		if _, ok := rsp["msg"].(string); !ok {
			return []byte(bod), fmt.Errorf("Helloworld is not saying hello, response body: '%v'", bod)
		}
		return []byte(bod), nil
	}, 90*time.Second); err != nil {
		return
	}
}

func TestServerAuth(t *testing.T) {
	TrySuite(t, ServerAuth, retryCount)
}

func ServerAuth(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// Execute first command in read to wait for store service
	// to start up
	if err := Try("Calling micro auth list accounts", t, func() ([]byte, error) {
		outp, err := cmd.Exec("auth", "list", "accounts")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "admin") {
			return outp, fmt.Errorf("Output should contain default admin account")
		}
		return outp, nil
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Calling micro auth list rules", t, func() ([]byte, error) {
		outp, err := cmd.Exec("auth", "list", "rules")
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

	if err := Try("Try to get token with default account", t, func() ([]byte, error) {
		outp, err := cmd.Exec("call", "auth", "Auth.Token", `{"id":"admin","secret":"micro"}`)
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
	TrySuite(t, testServerLockdown, retryCount)
}

func testServerLockdown(t *T) {
	t.Parallel()
	serv := NewServer(t)
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	lockdownSuite(serv, t)
}

func lockdownSuite(serv Server, t *T) {
	cmd := serv.Command()

	// Execute first command in read to wait for store service
	// to start up
	ns, err := namespace.Get(serv.Env())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Namespace is", ns)

	_, rsp, _ := curl(serv, "micro", "store/list")
	if rsp == nil {
		t.Fatal(rsp, errors.New("store list should have response"))
	}
	if val, ok := rsp["Code"].(float64); !ok || val != 401 {
		t.Fatal(rsp, errors.New("store list should be closed"), val)
	}

	Login(serv, t, "admin", "micro")

	email := "me@email.com"
	pass := "mystrongpass"

	outp, err := cmd.Exec("auth", "create", "account", "--secret", pass, email)
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = cmd.Exec("auth", "create", "rule", "--access=granted", "--scope=*", "--resource=*:*:*", "onlyloggedin")
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = cmd.Exec("auth", "create", "rule", "--access=granted", "--resource=service:auth:Auth.Token", "authpublic")
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = cmd.Exec("auth", "delete", "rule", "admin")
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	// set the local config file to be the same as the one micro will be configured to use.
	// todo: consider adding a micro logout command.
	config.SetConfig(cmd.Config)
	outp, err = cmd.Exec("logout")
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	if err := Try("Listing rules should fail before login", t, func() ([]byte, error) {
		outp, err := cmd.Exec("auth", "list", "rules")
		if err == nil {
			return outp, errors.New("List rules should fail")
		}
		return outp, nil
	}, 40*time.Second); err != nil {
		return
	}

	// auth rules are cached so this could take a few seconds (until the authpublic rule takes
	// effect in both the proxy and the auth service)
	if err := Try("Logging in with "+email, t, func() ([]byte, error) {
		out, err := cmd.Exec("login", "--email", email, "--password", pass)
		if err != nil {
			return out, err
		}
		if !strings.Contains(string(out), "Success") {
			return out, errors.New("Login output does not contain 'Success'")
		}
		return out, err
	}, 40*time.Second); err != nil {
		return
	}

	if err := Try("Listing rules should pass after login", t, func() ([]byte, error) {
		outp, err := cmd.Exec("auth", "list", "rules")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "onlyloggedin") || !strings.Contains(string(outp), "authpublic") {
			return outp, errors.New("Can't find rules")
		}
		return outp, err
	}, 40*time.Second); err != nil {
		return
	}
}

func TestPasswordChange(t *testing.T) {
	TrySuite(t, changePassword, retryCount)
}

func changePassword(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	newPass := "shinyNewPass"

	// Bad password should not succeed
	outp, err := cmd.Exec("user", "set", "password", "--old-password", "micro121212", "--new-password", newPass)
	if err == nil {
		t.Fatal("Incorrect existing password should make password change fail")
		return
	}

	outp, err = cmd.Exec("user", "set", "password", "--old-password", "micro", "--new-password", newPass)
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	time.Sleep(3 * time.Second)
	outp, err = cmd.Exec("login", "--email", "admin", "--password", "micro")
	if err == nil {
		t.Fatal("Old password should not be usable anymore")
		return
	}
	outp, err = cmd.Exec("login", "--email", "admin", "--password", newPass)
	if err != nil {
		t.Fatal(string(outp))
		return
	}
}

// TestUsernameLogin tests whether we can login using both ID and username e.g. UUID and email
func TestUsernameLogin(t *testing.T) {
	TrySuite(t, testUsernameLogin, retryCount)
}

func testUsernameLogin(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	outp, err := cmd.Exec("call", "auth", "Auth.Generate", `{"id":"someID", "name":"someUsername", "secret":"password"}`)
	if err != nil {
		t.Fatalf("Error generating account %s %s", string(outp), err)
	}
	outp, err = cmd.Exec("login", "--username", "someUsername", "--password", "password")
	if err != nil {
		t.Fatalf("Error logging in with user name %s %s", string(outp), err)
	}
	outp, err = cmd.Exec("login", "--username", "someID", "--password", "password")
	if err != nil {
		t.Fatalf("Error logging in with ID %s %s", string(outp), err)
	}
	// test the email alias
	outp, err = cmd.Exec("login", "--email", "someID", "--password", "password")
	if err != nil {
		t.Fatalf("Error logging in with ID %s %s", string(outp), err)
	}

	// test we can't create an account with the same name but different ID
	outp, err = cmd.Exec("call", "auth", "Auth.Generate", `{"id":"someID2", "name":"someUsername", "secret":"password1"}`)
	if err == nil {
		// shouldn't let us create something with the same username
		t.Fatalf("Expected error when generating account %s %s", string(outp), err)
	}

	outp, err = cmd.Exec("auth", "list", "accounts")
	if err != nil {
		t.Fatalf("Error listing accounts %s %s", string(outp), err)
	}
	if !strings.Contains(string(outp), "someUsername") {
		t.Fatalf("Error listing accounts, name is missing from %s", string(outp))
	}

	outp, err = cmd.Exec("login", "--username", "someID", "--password", "password")
	if err != nil {
		t.Fatalf("Error logging in with ID %s %s", string(outp), err)
	}

	// make sure user sees username and not ID
	outp, err = cmd.Exec("user")
	if err != nil {
		t.Fatalf("Error running user command %s %s", string(outp), err)
	}
	if !strings.Contains(string(outp), "someUsername") {
		t.Fatalf("Error running user command. Unexpected result %s", string(outp))
	}
	// make sure user sees username and not ID
	outp, err = cmd.Exec("user", "config")
	if err != nil {
		t.Fatalf("Error running user config command %s %s", string(outp), err)
	}
	if !strings.Contains(string(outp), "someUsername") {
		t.Fatalf("Error running user config command. Unexpected result %s", string(outp))
	}
	// make sure change password works correctly for username
	outp, err = cmd.Exec("user", "set", "password", "--old-password", "password", "--new-password", "password1")
	if err != nil {
		t.Fatalf("Error changing password %s %s", string(outp), err)
	}

	outp, err = cmd.Exec("login", "--username", "someUsername", "--password", "password1")
	if err != nil {
		t.Fatalf("Error changing password %s %s", string(outp), err)
	}

	outp, err = cmd.Exec("run", "github.com/micro/examples/helloworld")
	if err != nil {
		t.Fatalf("Error running helloworld %s %s", string(outp), err)
	}
	Try("Check helloworld status", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, fmt.Errorf("Error getting status %s", err)
		}
		if !strings.Contains(string(outp), "owner=someUsername") {
			return outp, fmt.Errorf("Can't find owner")
		}
		return nil, nil
	}, 30*time.Second)

}
