// +build integration

package test

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/token"
)

// Test for making sure config and store values across namespaces
// are correctly isolated
func TestNamespaceConfigIsolation(t *testing.T) {
	TrySuite(t, testNamespaceConfigIsolation, retryCount)
}

func testNamespaceConfigIsolation(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	testNamespaceConfigIsolationSuite(serv, t)
}

func testNamespaceConfigIsolationSuite(serv Server, t *T) {
	err := namespace.Add(serv.EnvName(), serv.EnvName())
	if err != nil {
		t.Fatal(err)
		return
	}
	err = namespace.Set(serv.EnvName(), serv.EnvName())
	if err != nil {
		t.Fatal(err)
		return
	}

	Login(serv, t, "default", "password")
	if t.failed {
		return
	}

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

	err = namespace.Add("random", serv.EnvName())
	if err != nil {
		t.Fatal(err)
		return
	}
	err = namespace.Set("random", serv.EnvName())
	if err != nil {
		t.Fatal(err)
		return
	}
	err = token.Remove(serv.EnvName())
	if err != nil {
		t.Fatal(err)
		return
	}

	// This call is only here to trigger default account generation
	exec.Command("micro", serv.EnvFlag(), "auth", "list", "accounts").CombinedOutput()

	Login(serv, t, "default", "password")
	if t.failed {
		return
	}

	if err := Try("reading 'somekey' should not be found with this account", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", serv.EnvFlag(), "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err == nil {
			return outp, errors.New("getting somekey should fail")
		}
		if !strings.Contains(string(outp), "Not found") {
			return outp, errors.New("Expected 'not found\n'")
		}
		return outp, nil
	}, 8 * time.Second); err != nil {
		return
	}

	// Log back to original namespace and see if value is already there

	// orignal namespace matchesthe env name
	err = namespace.Set(serv.EnvName(), serv.EnvName())
	if err != nil {
		t.Fatal(err)
		return
	}
	err = token.Remove(serv.EnvName())
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := Login(serv, t, "default", "password"); err != nil {
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
}
