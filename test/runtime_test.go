// +build integration

package test

import (
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

type server struct {
	cmd *exec.Cmd
	t   *testing.T
}

func newServer(t *testing.T) server {
	outp, err := exec.Command("rm", "-rf", "/tmp/micro/store").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}
	return server{cmd: exec.Command("micro", "server"), t: t}
}

func (s server) launch() {
	go func() {
		if err := s.cmd.Start(); err != nil {
			s.t.Fatal(err)
		}
	}()
	time.Sleep(1300 * time.Millisecond)
}

func (s server) close() {
	if s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGTERM)
	}
}

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

	serv := newServer(t)
	serv.launch()
	defer serv.close()

	outp, err = exec.Command("micro", "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
	if err != nil {
		t.Fatalf("Call to runtime read should succeed, err: %v, outp: %v", err, string(outp))
	}
}

func TestMicroRun(t *testing.T) {
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	psCmd := exec.Command("micro", "ps")
	outp, err = psCmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	// The started service should have the runtime name of "test/example-service",
	// as the runtime name is the relative path inside a repo.
	if !strings.Contains(string(outp), "test/example-service") {
		t.Fatal(string(outp))
	}
}
