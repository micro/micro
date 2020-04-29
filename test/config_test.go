// +build integration

package test

import (
	"errors"
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	try("Calling micro store read", t, func() ([]byte, error) {
		getCmd := exec.Command("micro", "config", "get", "somekey")
		outp, err := getCmd.CombinedOutput()
		if err == nil {
			return outp, errors.New("config gete should fail")
		}
		if string(outp) != "not found\n" {
			return outp, fmt.Errorf("Output should be 'not found\n', got %v", string(outp))
		}
		return outp, nil
	}, 5*time.Second)

	setCmd := exec.Command("micro", "config", "set", "somekey", "val1")
	outp, err := setCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
	}

	getCmd := exec.Command("micro", "config", "get", "somekey")
	outp, err = getCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(outp) != "val1\n" {
		t.Fatalf("Expected 'val1\n', got: '%v'", string(outp))
	}

	delCmd := exec.Command("micro", "config", "del", "somekey")
	outp, err = delCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(outp) != "" {
		t.Fatalf("Expected '', got: '%v'", string(outp))
	}

	getCmd = exec.Command("micro", "config", "get", "somekey")
	outp, err = getCmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Config get should fail: %v", string(outp))
	}
	if string(outp) != "not found\n" {
		t.Fatalf("Expected 'not found\n', got: '%v'", string(outp))
	}

	// Testing dot notation
	setCmd = exec.Command("micro", "config", "set", "someotherkey.subkey", "otherval1")
	outp, err = setCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(outp) != "" {
		t.Fatalf("Expected no output, got: %v", string(outp))
	}

	getCmd = exec.Command("micro", "config", "get", "someotherkey.subkey")
	outp, err = getCmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}
	if string(outp) != "otherval1\n" {
		t.Fatalf("Expected 'otherval1\n', got: '%v'", string(outp))
	}
}
