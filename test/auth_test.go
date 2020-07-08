// +build integration

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/token"
)

func TestServerAuth(t *testing.T) {
	trySuite(t, testServerAuth, retryCount)
}

func testServerAuth(t *t) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

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
	serv.launch()
	defer serv.close()

	basicAuthSuite(serv, t)
}

func basicAuthSuite(serv server, t *t) {
	// Execute first command in read to wait for store service
	// to start up
	try("Calling micro auth list accounts", t, func() ([]byte, error) {
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
	}, 15*time.Second)

	try("Calling micro auth list rules", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "auth", "list", "rules")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "default") {
			return outp, fmt.Errorf("Output should contain default rule")
		}
		return outp, nil
	}, 8*time.Second)

	try("Try to get token with default account", t, func() ([]byte, error) {
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
	}, 8*time.Second)
}

// Test for making sure config and store values across namespaces
// are correctly isolated
func TestNamespaceIsolation(t *testing.T) {
	trySuite(t, testNamespaceIsolation, retryCount)
}

func testNamespaceIsolation(t *t) {
	t.Parallel()
	serv := newServer(t, options{
		auth: "jwt",
	})
	serv.launch()
	defer serv.close()

	testNamespaceIsolationSuite(serv, t)
}

func testNamespaceIsolationSuite(serv server, t *t) {
	// Execute first command in read to wait for store service
	// to start up
	try("Calling micro auth list accounts", t, func() ([]byte, error) {
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
	}, 15*time.Second)

	login(serv, t, "default", "password")
	if t.failed {
		return
	}

	try("Calling micro config set", t, func() ([]byte, error) {
		setCmd := exec.Command("micro", serv.envFlag(), "config", "set", "somekey", "val1")
		outp, err := setCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		return outp, err
	}, 5*time.Second)

	try("micro config get somekey", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.envFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "val1\n" {
			return outp, errors.New("Expected 'val1\n'")
		}
		return outp, err
	}, 8*time.Second)

	namespace.Set(serv.envName, "random")

	try("Calling micro auth list accounts", t, func() ([]byte, error) {
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
	}, 15*time.Second)

	try("reading 'somekey' should not be found with this account", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.envFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err == nil {
			return outp, errors.New("getting somekey should fail")
		}
		if string(outp) != "not found\n" {
			return outp, errors.New("Expected 'not found\n'")
		}
		return outp, nil
	}, 8*time.Second)
}

func login(serv server, t *t, email, password string) {
	try("Logging in", t, func() ([]byte, error) {
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
