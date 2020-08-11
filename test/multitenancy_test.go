// +build integration

package test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
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
	err := namespace.Add(serv.Env(), serv.Env())
	if err != nil {
		t.Fatal(err)
		return
	}
	err = namespace.Set(serv.Env(), serv.Env())
	if err != nil {
		t.Fatal(err)
		return
	}

	Login(serv, t, "admin", "micro")
	if t.failed {
		return
	}

	cmd := serv.Command()

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
	currNamespace, _ := cmd.Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err := ChangeNamespace(cmd, serv.Env(), "random"); err != nil {
		t.Fatalf("Error changing namespace %s", err)
	}

	// This call is only here to trigger default account generation
	cmd.Exec("auth", "list", "accounts")

	Login(serv, t, "admin", "micro")
	if t.failed {
		return
	}

	if err := Try("reading 'somekey' should not be found with this account", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "get", "somekey")
		if err == nil {
			return outp, errors.New("getting somekey should fail")
		}
		if !strings.Contains(string(outp), "Unauthorized") {
			return outp, errors.New("Expected 'Unauthorized\n'")
		}
		return outp, nil
	}, 8*time.Second); err != nil {
		return
	}

	// Log back to original namespace and see if value is already there
	if err := ChangeNamespace(cmd, serv.Env(), string(currNamespace)); err != nil {
		t.Fatalf("Error changing namespace %s", err)
	}

	if err := Login(serv, t, "admin", "micro"); err != nil {
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
}
