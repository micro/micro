// +build integration

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/token"
)

// Test for making sure config and store values across namespaces
// are correctly isolated
func TestNamespaceConfigIsolation(t *testing.T) {
	trySuite(t, testNamespaceConfigIsolation, retryCount)
}

func testNamespaceConfigIsolation(t *t) {
	t.Parallel()
	serv := newServer(t)
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	testNamespaceConfigIsolationSuite(serv, t)
}

func testNamespaceConfigIsolationSuite(serv server, t *t) {
	err := namespace.Add(serv.envName, serv.envName)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = namespace.Set(serv.envName, serv.envName)
	if err != nil {
		t.Fatal(err)
		return
	}

	// This call is only here to trigger default account generation
	exec.Command("micro", serv.envFlag(), "auth", "list", "accounts").CombinedOutput()

	login(serv, t, "default", "password")
	if t.failed {
		return
	}

	if err := try("Calling micro config set", t, func() ([]byte, error) {
		setCmd := exec.Command("micro", serv.envFlag(), "config", "set", "somekey", "val1")
		outp, err := setCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		return outp, err
	}, 5*time.Second); err != nil {
		return
	}

	if err := try("micro config get somekey", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.envFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "val1\n" {
			return outp, errors.New("Expected 'val1\n'")
		}
		return outp, err
	}, 8*time.Second); err != nil {
		return
	}

	err = namespace.Add("random", serv.envName)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = namespace.Set("random", serv.envName)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = token.Remove(serv.envName)
	if err != nil {
		t.Fatal(err)
		return
	}

	// This call is only here to trigger default account generation
	exec.Command("micro", serv.envFlag(), "auth", "list", "accounts").CombinedOutput()

	login(serv, t, "default", "password")
	if t.failed {
		return
	}

	if err := try("reading 'somekey' should not be found with this account", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.envFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err == nil {
			return outp, errors.New("getting somekey should fail")
		}
		if string(outp) != "not found\n" {
			return outp, errors.New("Expected 'not found\n'")
		}
		return outp, nil
	}, 8*time.Second); err != nil {
		return
	}

	// Log back to original namespace and see if value is already there

	// orignal namespace matchesthe env name
	err = namespace.Set(serv.envName, serv.envName)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = token.Remove(serv.envName)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := login(serv, t, "default", "password"); err != nil {
		return
	}

	if err := try("micro config get somekey", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.envFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "val1\n" {
			return outp, errors.New("Expected 'val1\n'")
		}
		return outp, err
	}, 8*time.Second); err != nil {
		return
	}
}

func login(serv server, t *t, email, password string) error {
	return try("Logging in", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.auth", "Auth.Token", fmt.Sprintf(`{"id":"%v","secret":"%v"}`, email, password))
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]interface{}{}
		err = json.Unmarshal(outp, &rsp)
		tok, ok := rsp["token"].(map[string]interface{})
		if !ok {
			return outp, errors.New("Can't find token")
		}
		if _, ok = tok["access_token"].(string); !ok {
			return outp, fmt.Errorf("Can't find access token")
		}
		if _, ok = tok["refresh_token"].(string); !ok {
			return outp, fmt.Errorf("Can't find access token")
		}
		if _, ok = tok["refresh_token"].(string); !ok {
			return outp, fmt.Errorf("Can't find refresh token")
		}
		if _, ok = tok["expiry"].(string); !ok {
			return outp, fmt.Errorf("Can't find access token")
		}
		exp, err := strconv.ParseInt(tok["expiry"].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		token.Save(serv.envName, &auth.Token{
			AccessToken:  tok["access_token"].(string),
			RefreshToken: tok["refresh_token"].(string),
			Expiry:       time.Unix(exp, 0),
		})
		return outp, nil
	}, 8*time.Second)
}
