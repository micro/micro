// +build integration

package test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	TrySuite(t, testStore, 5)
}

func testStore(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// Execute first command in read to wait for store service
	// to start up
	if err := Try("Calling micro store read", t, func() ([]byte, error) {
		outp, err := cmd.Exec("store", "read", "somekey")
		if err == nil {
			return outp, errors.New("store read should fail")
		}
		if !strings.Contains(string(outp), "Not found") {
			return outp, fmt.Errorf("Output should be 'Not found', got %v", string(outp))
		}
		return outp, nil
	}, 8*time.Second); err != nil {
		return
	}

	outp, err := cmd.Exec("store", "write", "somekey", "val1")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	outp, err = cmd.Exec("store", "read", "somekey")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "val1\n" {
		t.Fatalf("Expected 'val1\n', got: '%v'", string(outp))
		return
	}

	outp, err = cmd.Exec("store", "delete", "somekey")
	if err != nil {
		t.Fatal(err)
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected '', got: '%v'", string(outp))
		return
	}

	outp, err = cmd.Exec("store", "read", "somekey")
	if err == nil {
		t.Fatalf("store read should fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "Not found") {
		t.Fatalf("Expected 'Not found\n', got: '%v'", string(outp))
		return
	}

	// Test prefixes
	outp, err = cmd.Exec("store", "write", "somekey1", "val1")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	outp, err = cmd.Exec("store", "write", "somekey2", "val2")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
		return
	}

	// Read exact key
	outp, err = cmd.Exec("store", "read", "somekey")
	if err == nil {
		t.Fatalf("store read should fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "Not found") {
		t.Fatalf("Expected 'Not found\n', got: '%v'", string(outp))
		return
	}

	outp, err = cmd.Exec("store", "read", "--prefix", "somekey")
	if err != nil {
		t.Fatalf("store prefix read not should fail: %v", string(outp))
		return
	}
	if string(outp) != "val1\nval2\n" {
		t.Fatalf("Expected output not present, got: '%v'", string(outp))
		return
	}

	outp, err = cmd.Exec("store", "read", "-v", "--prefix", "somekey")
	if err != nil {
		t.Fatalf("store prefix read not should fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "somekey1") || !strings.Contains(string(outp), "somekey2") ||
		!strings.Contains(string(outp), "val1") || !strings.Contains(string(outp), "val2") {
		t.Fatalf("Expected output not present, got: '%v'", string(outp))
		return
	}

	outp, err = cmd.Exec("store", "list")
	if err != nil {
		t.Fatalf("store list should not fail: %v", string(outp))
		return
	}
	if !strings.Contains(string(outp), "somekey1") || !strings.Contains(string(outp), "somekey2") {
		t.Fatalf("Expected output not present, got: '%v'", string(outp))
		return
	}

}

func TestStoreImpl(t *testing.T) {
	TrySuite(t, testStoreImpl, 3)
}

func testStoreImpl(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "./services/test/kv")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find store", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if !statusRunning("kv", "latest", outp) {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 90*time.Second); err != nil {
		return
	}

	if err := Try("Check logs", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "kv")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "Listening on") {
			return outp, fmt.Errorf("Service not ready")
		}
		return nil, nil
	}, 60*time.Second); err != nil {
		return
	}
	outp, err = cmd.Exec("call", "--request_timeout=15s", "example", "Example.TestExpiry")
	if err != nil {
		t.Fatalf("Error %s, %s", err, outp)
	}

	outp, err = cmd.Exec("call", "--request_timeout=15s", "example", "Example.TestList")
	if err != nil {
		t.Fatalf("Error %s, %s", err, outp)
	}

	outp, err = cmd.Exec("call", "--request_timeout=15s", "example", "Example.TestListLimit")
	if err != nil {
		t.Fatalf("Error %s, %s", err, outp)
	}
	outp, err = cmd.Exec("call", "--request_timeout=15s", "example", "Example.TestListOffset")
	if err != nil {
		t.Fatalf("Error %s, %s", err, outp)
	}
}

func TestBlobStore(t *testing.T) {
	TrySuite(t, testBlobStore, retryCount)
}

func testBlobStore(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "./services/test/blob-store")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find blob-store", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("blob-store", "latest", outp) {
			return outp, errors.New("Can't find blob-store service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Check logs", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "blob-store")
		if err != nil {
			return nil, err
		}
		if !strings.Contains(string(outp), "Read from blob store: world") {
			return outp, fmt.Errorf("Didn't read from the blob store")
		}
		return nil, nil
	}, 60*time.Second); err != nil {
		return
	}
}
