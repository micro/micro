// +build kind

package test

import (
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/token"
)

func init() {
	testFilter = []string{
		"TestRunGithubSource",
		"TestStore",
		"TestCorruptedTokenLogin",
	}
	maxTimeMultiplier = 3
	isParallel = false // in theory should work in parallel
	newSrv = newK8sServer
	retryCount = 1
}

func newK8sServer(t *t, fname string, opts ...options) testServer {
	portnum := rand.Intn(maxPort-minPort) + minPort

	s := &testK8sServer{testServerBase{
		t:       t,
		envNm:   strings.ToLower(fname),
		portNum: portnum,
		cmd:     exec.Command("kubectl", "port-forward", "--namespace", "default", "svc/micro-proxy", fmt.Sprintf("%d:8081", portnum)),
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

	// login to admin account
	if err := login(s, s.t, "default", "password"); err != nil {
		s.t.Fatalf("Error logging in %s", err)
		return err
	}

	if err := try("Calling micro server", s.t, func() ([]byte, error) {
		outp, err := exec.Command("micro", s.envFlag(), "services").CombinedOutput()
		if !strings.Contains(string(outp), "runtime") ||
			!strings.Contains(string(outp), "registry") ||
			!strings.Contains(string(outp), "broker") ||
			!strings.Contains(string(outp), "config") ||
			!strings.Contains(string(outp), "debug") ||
			!strings.Contains(string(outp), "proxy") ||
			!strings.Contains(string(outp), "auth") ||
			!strings.Contains(string(outp), "store") {
			return outp, errors.New("Not ready")
		}

		return outp, err
	}, 60*time.Second); err != nil {
		return err
	}

	// generate a new admin account for the env : user=ENV_NAME pass=password
	req := fmt.Sprintf(`{"id":"%s", "secret":"password", "options":{"namespace":"%s"}}`, s.envName(), s.namespace)
	outp, err := exec.Command("micro", s.envFlag(), "call", "go.micro.auth", "Auth.Generate", req).CombinedOutput()
	if err != nil && !strings.Contains(string(outp), "already exists") { // until auth.Delete is implemented
		s.t.Fatalf("Error generating auth: %s, %s", err, outp)
		return err
	}

	// remove the admin token
	token.Remove(s.envName())

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
