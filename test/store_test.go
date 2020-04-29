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
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	// Execute first command in read to wait for store service
	// to start up
	try("Calling micro store read", t, func() ([]byte, error) {
		readCmd := exec.Command("micro", "store", "read", "somekey")
		outp, err := readCmd.CombinedOutput()
		if err == nil {
			return outp, errors.New("store read should fail")
		}
		if string(outp) != "not found\n" {
			return outp, fmt.Errorf("Output should be 'not found\n', got %v", string(outp))
		}
		return outp, nil
	}, 5*time.Second)

	writeCmd := exec.Command("micro", "store", "write", "somekey", "val1")
	outp, err := writeCmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
	}

	readCmd := exec.Command("micro", "store", "read", "somekey")
	outp, err = readCmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}
	if !strings.Contains(string(outp), "val1") {
		t.Fatalf("Expected 'val1\n', got: '%v'", string(outp))
	}

	delCmd := exec.Command("micro", "store", "delete", "somekey")
	outp, err = delCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(outp) != "" {
		t.Fatalf("Expected '', got: '%v'", string(outp))
	}

	readCmd = exec.Command("micro", "store", "read", "somekey")
	outp, err = readCmd.CombinedOutput()
	if err == nil {
		t.Fatalf("store read should fail: %v", string(outp))
	}
	if string(outp) != "not found\n" {
		t.Fatalf("Expected 'not found\n', got: '%v'", string(outp))
	}

}
