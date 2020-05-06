// +build integration

package test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestServerAuth(t *testing.T) {
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	// Execute first command in read to wait for store service
	// to start up
	try("Calling micro auth list accounts", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", "auth", "list", "accounts")
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
		readCmd := exec.Command("micro", "auth", "list", "rules")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "*:*:*:*") {
			return outp, fmt.Errorf("Output should contain default rule")
		}
		return outp, nil
	}, 8*time.Second)

	try("Calling micro auth list rules", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", "auth", "list", "rules")
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "*:*:*:*") {
			return outp, fmt.Errorf("Output should contain default rule")
		}
		return outp, nil
	}, 8*time.Second)

	runCmd := exec.Command("micro", "run", "helloworld")
	_, err := runCmd.CombinedOutput()
	if err != nil {
		return t.Fatal(err)
	}

	try("Call hello world", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", "call", "go.micro.service.helloworld", "Helloworld.Call", `{"name": "Joe"}`)
		outp, err := callCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hello Joe" {
			return outp, errors.New("Helloworld resonse is unexpected")
		}
		return outp, err
	}, 20*time.Second)
}
