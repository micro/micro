// +build kind

package test

import (
	"errors"
	"fmt"
	"math/rand"
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
		"TestPublicAPI",
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

	s := &testK8sServer{ServerBase{
		dir:       filepath.Dir(configFile),
		config:    configFile,
		t:         t,
		env:       options.Namespace,
		proxyPort: portnum,
		opts:      options,
		cmd:       exec.Command("kubectl", "port-forward", "--namespace", "default", "svc/micro-proxy", fmt.Sprintf("%d:443", portnum)),
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
			!strings.Contains(string(outp), "updater") ||
			!strings.Contains(string(outp), "store") {
			return outp, errors.New("Not ready")
		}

		return outp, err
	}, 60*time.Second); err != nil {
		return err
	}

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

func TestPublicAPI(t *testing.T) {
	TrySuite(t, testPublicAPI, retryCount)
}

func testPublicAPI(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	err := ChangeNamespace(cmd, serv.Env(), "random-namespace")
	if err != nil {
		t.Fatal(err)
		return
	}
	// login to admin account
	if err := Login(serv, t, "admin", "micro"); err != nil {
		t.Fatalf("Error logging in %s", err)
		return
	}

	outp, err := cmd.Exec("user", "namespace")
	if err != nil || strings.TrimSpace(string(outp)) != "random-namespace" {
		t.Fatal(string(outp), err)
		return
	}

	if err := Try("Find helloworld", t, func() ([]byte, error) {
		outp, err = cmd.Exec("user")
		if err != nil || strings.TrimSpace(string(outp)) != "admin" {
			return outp, err
		}
		return outp, err
	}, 10*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld", t, func() ([]byte, error) {
		return cmd.Exec("run", "helloworld")
	}, 5*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if !statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("Can't find example helloworld in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Call helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("helloworld", "--name=joe")
		if err != nil {
			outp1, _ := cmd.Exec("logs", "helloworld")
			return append(outp, outp1...), err
		}
		if !strings.Contains(string(outp), "Msg") {
			return outp, err
		}
		return outp, err
	}, 90*time.Second); err != nil {
		return
	}

	if err := Try("curl helloworld", t, func() ([]byte, error) {
		bod, rsp, err := curl(serv, "random-namespace", "helloworld?name=Jane")
		if rsp == nil {
			return []byte(bod), fmt.Errorf("helloworld should have response, err: %v", err)
		}
		if _, ok := rsp["Msg"].(string); !ok {
			return []byte(bod), fmt.Errorf("Helloworld is not saying hello, response body: '%v'", bod)
		}
		return []byte(bod), nil
	}, 90*time.Second); err != nil {
		return
	}
}
