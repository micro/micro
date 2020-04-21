// +build integration

package test

import (
	"log"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestMicroServerModeCall(t *testing.T) {
	outp, err := exec.Command("micro", "env", "set", "server").CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to set env to server, err: %v, output: %v", err, string(outp))
	}

	callCmd := exec.Command("micro", "call", "go.micro.runtime", "Runtime.Read", "{}")
	outp, err = callCmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Call to server should fail, got no error, output: %v", string(outp))
	}

	serverCmd := exec.Command("micro", "server")
	go func() {
		if err := serverCmd.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	defer func() {
		if serverCmd.Process != nil {
			serverCmd.Process.Signal(syscall.SIGTERM)
		}
	}()
	time.Sleep(4 * time.Second)

	outp, err = exec.Command("micro", "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
	if err != nil {
		t.Fatalf("Call to runtime read should succeed, err: %v, outp: %v", err, string(outp))
	}
}
