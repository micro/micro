// +build kind

package test

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func init() {
	testFilter = []string{
		"TestRunGithubSource",
		"TestStore",
		// @todo Reactivate this once source to running works in kind
		// "TestStoreImpl",
		"TestCorruptedTokenLogin",
		"TestRunPrivateSource",
		"TestEventsStream",
		"TestRPC",
	}
	maxTimeMultiplier = 3
	isParallel = false // in theory should work in parallel
	newSrv = newK8sServer
	retryCount = 1
}

func newK8sServer(t *T, fname string, opts ...Option) Server {
	options := Options{
		Namespace: strings.ToLower(fname),
		Login:     false,
	}
	for _, o := range opts {
		o(&options)
	}

	portnum := rand.Intn(maxPort-minPort) + minPort
	configFile := configFile(fname)

	var cmd *exec.Cmd
	if v := os.Getenv("IN_HELM_TEST"); len(v) > 0 {
		cmd = exec.Command("kubectl", "port-forward", "--namespace", "micro", "svc/proxy", fmt.Sprintf("%d:443", portnum))
	} else {
		cmd = exec.Command("kubectl", "port-forward", "--namespace", "default", "svc/micro-proxy", fmt.Sprintf("%d:8081", portnum))
	}

	s := &testK8sServer{ServerBase{
		dir:       filepath.Dir(configFile),
		config:    configFile,
		t:         t,
		env:       options.Namespace,
		proxyPort: portnum,
		opts:      options,
		cmd:       cmd,
	}}
	s.namespace = s.env

	return s
}

type testK8sServer struct {
	ServerBase
}

func (s *testK8sServer) Run() error {
	if err := s.ServerBase.Run(); err != nil {
		return err
	}

	ChangeNamespace(s.Command(), s.Env(), "micro")

	// login to admin account
	if err := Login(s, s.t, "admin", "micro"); err != nil {
		s.t.Fatalf("Error logging in %s", err)
		return err
	}

	if err := Try("Calling micro server", s.t, func() ([]byte, error) {
		outp, err := s.Command().Exec("services")
		if !strings.Contains(string(outp), "runtime") ||
			!strings.Contains(string(outp), "registry") ||
			!strings.Contains(string(outp), "broker") ||
			!strings.Contains(string(outp), "config") ||
			!strings.Contains(string(outp), "proxy") ||
			!strings.Contains(string(outp), "auth") ||
			!strings.Contains(string(outp), "store") {
			return outp, errors.New("Not ready")
		}

		return outp, err
	}, 60*time.Second); err != nil {
		return err
	}

	// generate account in new namespace
	// ignore errors because it is not an idempotent call
	s.Command().Exec("auth", "create", "account", "--secret", "micro", "--namespace", s.Env(), "--scopes", "admin", "admin")

	// switch to the namespace
	ChangeNamespace(s.Command(), s.Env(), s.Env())

	// login to the admin account which is generated for each namespace
	if s.opts.Login {
		Login(s, s.t, "admin", "micro")
	}

	return nil
}

func (s *testK8sServer) Close() {
	s.ServerBase.Close()
	// kill the port forward
	s.cmd.Process.Kill()
}

func TestDeleteOwnAccount(t *testing.T) {
	TrySuite(t, testDeleteOwnAccount, retryCount)
}

func testDeleteOwnAccount(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	outp, err := cmd.Exec("auth", "delete", "account", "admin")
	if err == nil {
		t.Fatal(string(outp))
	}
}
