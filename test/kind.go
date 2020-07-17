// +build kind

package test

import (
	"fmt"
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
	s.loginOpts = fmt.Sprintf(`{"namespace":"%s"}`, s.envName())

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
	// login to admin account
	login(s, t, "default", "password", loginOptions{admin: true})

	// generate a new admin account for the env : user=ENV_NAME pass=password
	req := fmt.Sprintf(`{"id":"%s", "secret":"password", "options":{"namespace":"%s"}}`, s.envName(), s.envName())
	outp, err := exec.Command("micro", s.envFlag(), "call", "go.micro.auth", "Auth.Generate", req).CombinedOutput()
	if err != nil && !strings.Contains(string(outp), "already exists") { // until auth.Delete is implemented
		t.Fatalf("Error generating auth: %s, %s", err, outp)
		return err
	}

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
	// tear down the services that were made
	// kubectl truncate namespace or somethin
	// kill the port forward
	s.cmd.Process.Kill()
	// exec.Command("kubectl", "delete", "namespace", s.envName)

}