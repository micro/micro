// +build integration

package test

import (
	"errors"
	"fmt"
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

	cmd := serv.Command()

	if err := Try("Calling micro config read", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "get", "somekey")
		if err == nil {
			return outp, errors.New("config gete should fail")
		}
		if !strings.Contains(string(outp), "Not found") {
			return outp, fmt.Errorf("Output should be 'Not found\n', got %v", string(outp))
		}
		return outp, nil
	}, 5*time.Second); err != nil {
		return
	}

	// This needs to be retried to the the "error listing rules"
	// error log output that happens when the auth service is not yet available.

	if err := Try("Calling micro config set", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "set", "somekey", "val1")
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

	if err := Try("micro config get somekey", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "get", "somekey")
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

	outp, err := cmd.Exec("config", "del", "somekey")
	if err != nil {
		t.Fatalf(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected '', got: '%v'", string(outp))
		return
	}

	if err := Try("micro config get somekey", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "get", "somekey")
		if err == nil {
			return outp, errors.New("getting somekey should fail")
		}
		if !strings.Contains(string(outp), "Not found") {
			return outp, errors.New("Expected 'Not found'")
		}
		return outp, nil
	}, 8*time.Second); err != nil {
		return
	}

	// Testing dot notation
	outp, err = cmd.Exec("config", "set", "someotherkey.subkey", "otherval1")
	if err != nil {
		t.Fatal(err)
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	if err := Try("micro config get someotherkey.subkey", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "get", "someotherkey.subkey")
		if err != nil {
			return outp, err
		}
		if string(outp) != "otherval1\n" {
			return outp, errors.New("Expected 'otherval1\n'")
		}
		return outp, err
	}, 8*time.Second); err != nil {
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

	cmd := serv.Command()

	// This needs to be retried to the the "error listing rules"
	// error log output that happens when the auth service is not yet available.
	if err := Try("Calling micro config set", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "set", "key.subkey", "val1")
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		outp, err = cmd.Exec("config", "set", "--secret", "key.subkey1", "42")
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		// Testing JSON escape of "<" etc chars.
		outp, err = cmd.Exec("config", "set", "--secret", "key.subkey2", "\"Micro Team <support@m3o.com>\"")
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		// Setting an other key for `val.String` test
		outp, err = cmd.Exec("config", "set", "key.subkey3", "\"Micro Test <test@m3o.com>\"")
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

	// check the value set correctly
	if err := Try("Calling micro config get", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "get", "key")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "val1") {
			return outp, fmt.Errorf("Expected output to contain val1, got: %v", string(outp))
		}

		return outp, err
	}, 5*time.Second); err != nil {
		return
	}

	outp, err := cmd.Exec("run", "./service/config-example")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Try logs read", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "config-example")
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "val1") {
			return outp, fmt.Errorf("Expected val1 in output, got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "42") {
			return outp, fmt.Errorf("Expected output to contain 42, got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "Micro Team <support@m3o.com>") {
			return outp, fmt.Errorf("Expected output to contain \"Micro Team <support@m3o.com>\", got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "Micro Test <test@m3o.com>") {
			return outp, fmt.Errorf("Expected output to contain \"Micro Test <test@m3o.com>\", got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "Default Hello") {
			return outp, fmt.Errorf("Expected output to contain \"Default Hello\", got: %v", string(outp))
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}
}
