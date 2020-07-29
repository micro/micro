// +build integration

package test

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	TrySuite(t, testConfig, retryCount)
}

func testConfig(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	if err := Try("Calling micro config read", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.EnvFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err == nil {
			return outp, errors.New("config gete should fail")
		}
		if !strings.Contains(string(outp), "Not found") {
			return outp, fmt.Errorf("Output should be 'not found\n', got %v", string(outp))
		}
		return outp, nil
	}, 5 * time.Second); err != nil {
		return
	}

	// This needs to be retried to the the "error listing rules"
	// error log output that happens when the auth service is not yet available.

	if err := Try("Calling micro config set", t, func() ([]byte, error) {
		setCmd := exec.Command("micro", serv.EnvFlag(), "config", "set", "somekey", "val1")
		outp, err := setCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		return outp, err
	}, 5 * time.Second); err != nil {
		return
	}

	if err := Try("micro config get somekey", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.EnvFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "val1\n" {
			return outp, errors.New("Expected 'val1\n'")
		}
		return outp, err
	}, 8 * time.Second); err != nil {
		return
	}

	delCmd := exec.Command("micro", serv.EnvFlag(), "config", "del", "somekey")
	outp, err := delCmd.CombinedOutput()
	if err != nil {
		t.Fatalf(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected '', got: '%v'", string(outp))
		return
	}

	if err := Try("micro config get somekey", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.EnvFlag(), "config", "get", "somekey")
		outp, err = getCmd.CombinedOutput()
		if err == nil {
			return outp, errors.New("getting somekey should fail")
		}
		if !strings.Contains(string(outp), "not found") {
			return outp, errors.New("Expected 'not found'")
		}
		return outp, nil
	}, 8 * time.Second); err != nil {
		return
	}

	// Testing dot notation
	setCmd := exec.Command("micro", serv.EnvFlag(), "config", "set", "someotherkey.subkey", "otherval1")
	outp, err = setCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	if err := Try("micro config get someotherkey.subkey", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.EnvFlag(), "config", "get", "someotherkey.subkey")
		outp, err = getCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "otherval1\n" {
			return outp, errors.New("Expected 'otherval1\n'")
		}
		return outp, err
	}, 8 * time.Second); err != nil {
		return
	}
}

func TestConfigReadFromService(t *testing.T) {
	TrySuite(t, testConfigReadFromService, retryCount)
}

func testConfigReadFromService(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	// This needs to be retried to the the "error listing rules"
	// error log output that happens when the auth service is not yet available.
	if err := Try("Calling micro config set", t, func() ([]byte, error) {
		setCmd := exec.Command("micro", serv.EnvFlag(), "config", "set", "key.subkey", "val1")
		outp, err := setCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		return outp, err
	}, 5 * time.Second); err != nil {
		return
	}

	// check the value set correctly
	if err := Try("Calling micro config get", t, func() ([]byte, error) {
		setCmd := exec.Command("micro", serv.EnvFlag(), "config", "get", "key")
		outp, err := setCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "val1") {
			return outp, fmt.Errorf("Expected output to contain val1, got: %v", string(outp))
		}

		return outp, err
	}, 5 * time.Second); err != nil {
		return
	}

	runCmd := exec.Command("micro", serv.EnvFlag(), "run", "./service/config")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Try logs read", t, func() ([]byte, error) {
		setCmd := exec.Command("micro", serv.EnvFlag(), "logs", "test/service/config")
		outp, err := setCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "val1") {
			return outp, fmt.Errorf("Expected val1 in output, got: %v", string(outp))
		}
		return outp, err
	}, 20 * time.Second); err != nil {
		return
	}
}
