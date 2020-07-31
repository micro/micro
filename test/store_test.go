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

func TestStore(t *testing.T) {
	TrySuite(t, testStore, 5)
}

func testStore(t *T) {
	t.Parallel()
	serv := NewServer(t)
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	if err := Login(serv, t, "default", "password"); err != nil {
		t.Fatalf("Failed to login %s", err)
	}

	// Execute first command in read to wait for store service
	// to start up
	if err := Try("Calling micro store read", t, func() ([]byte, error) {
		outp, err := serv.Command().Exec("store", "read", "somekey")
		if err == nil {
			return outp, errors.New("store read should fail")
		}
		if !strings.Contains(string(outp), "not found") {
			return outp, fmt.Errorf("Output should be 'not found', got %v", string(outp))
		}
		return outp, nil
	}, 8*time.Second); err != nil {
		return
	}

	outp, err := serv.Command().Exec("store", "write", "somekey", "val1")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	outp, err = serv.Command().Exec("store", "read", "somekey")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "val1\n" {
		t.Fatalf("Expected 'val1\n', got: '%v'", string(outp))
		return
	}

	outp, err = serv.Command().Exec("store", "delete", "somekey")
	if err != nil {
		t.Fatal(err)
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected '', got: '%v'", string(outp))
		return
	}

	outp, err = serv.Command().Exec("store", "read", "somekey")
	if err == nil {
		t.Fatalf("store read should fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "not found") {
		t.Fatalf("Expected 'not found\n', got: '%v'", string(outp))
		return
	}

	// Test prefixes
	outp, err = serv.Command().Exec("store", "write", "somekey1", "val1")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	outp, err = serv.Command().Exec("store", "write", "somekey2", "val2")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	// Read exact key
	outp, err = serv.Command().Exec("store", "read", "somekey")
	if err == nil {
		t.Fatalf("store read should fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "not found") {
		t.Fatalf("Expected 'not found\n', got: '%v'", string(outp))
		return
	}

	outp, err = serv.Command().Exec("store", "read", "--prefix", "somekey")
	if err != nil {
		t.Fatalf("store prefix read not should fail: %v", string(outp))
		return
	}
	if string(outp) != "val1\nval2\n" {
		t.Fatalf("Expected output not present, got: '%v'", string(outp))
		return
	}

	outp, err = serv.Command().Exec("store", "read", "-v", "--prefix", "somekey")
	if err != nil {
		t.Fatalf("store prefix read not should fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "somekey1") || !strings.Contains(string(outp), "somekey2") ||
		!strings.Contains(string(outp), "val1") || !strings.Contains(string(outp), "val2") {
		t.Fatalf("Expected output not present, got: '%v'", string(outp))
		return
	}

	outp, err = serv.Command().Exec("store", "list")
	if err != nil {
		t.Fatalf("store list should not fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "somekey1") || !strings.Contains(string(outp), "somekey2") {
		t.Fatalf("Expected output not present, got: '%v'", string(outp))
		return
	}
}
