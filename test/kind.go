// +build kind

package test

import (
	"os/exec"
	"strings"

	"github.com/micro/micro/v2/client/cli/namespace"
)

func init() {
	testFilter = []string{
		"TestRunGithubSource",
		"TestStore",
	}
	maxTimeMultiplier = 3
	isParallel = false // in theory should work in parallel
	newSrv = newK8sServer
	retryCount = 1
}

func newK8sServer(t *t, fname string, opts ...options) testServer {
	s := &testK8sServer{testServerBase{
		t:       t,
		envNm:   strings.ToLower(fname),
		portNum: 8081,
		cmd:     exec.Command("kubectl", "port-forward", "--namespace", "default", "svc/micro-proxy", "8081:8081"),
	}}
	s.namespace = s.envNm

	return s
}

type testK8sServer struct {
	testServerBase
}

func (s *testK8sServer) launch() error {
	if err := s.testServerBase.launch(); err != nil {
		return err
	}
	t := s.t

	// setup .micro config for access
	if err := namespace.Add(s.envName(), s.envName()); err != nil {
		t.Fatalf("Failed to add current namespace: %s", err)
		return err
	}
	if err := namespace.Set(s.envName(), s.envName()); err != nil {
		t.Fatalf("Failed to set current namespace: %s", err)
		return err
	}

	return nil
}

func (s *testK8sServer) close() {
	s.testServerBase.close()
	// kill the port forward
	s.cmd.Process.Kill()
}
