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
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	// Execute first command in read to wait for store service
	// to start up
	try("Calling micro auth list accounts", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "auth", "list", "accounts")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "admin") {
			return outp, fmt.Errorf("Output should contain admin")
		}
		return outp, nil
	}, 8*time.Second)

	try("Calling micro auth list rules", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "auth", "list", "rules")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "*:*:*:*") {
			return outp, fmt.Errorf("Output should contain default rule")
		}
		return outp, nil
	}, 8*time.Second)

	try("Calling micro auth list accounts", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "auth", "list", "accounts")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "admin") {
			return outp, fmt.Errorf("Output should contain default admin account")
		}
		return outp, nil
	}, 8*time.Second)

	accessToken := ""
	try("Try to get token with default account", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.auth", "Auth.Token", `{"id":"admin","secret":"Password1"}`)
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
		if accessToken, ok = token["access_token"].(string); !ok {
			return outp, fmt.Errorf("Can't find access token")
		}
		return outp, nil
	}, 5*time.Second)

	try("Try to log in with token we got", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.envFlag(), "login", "--token", accessToken)
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		return outp, nil
	}, 5*time.Second)
}
