// +build integration

package test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func Testroutes(t *testing.T) {
	TrySuite(t, testroutes, retryCount)
}

func testroutes(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "./services/test/routes/routes-server")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find routes-server in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}
		if !statusRunning("routes-server", "latest", outp) {
			return outp, errors.New("Can't find routes-server in runtime")
		}
		return nil, nil
	}, 120*time.Second); err != nil {
		return
	}

	if err := Try("Find routes service in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "routes") {
			return outp, errors.New("Can't find routes service in registry")
		}
		return nil, nil
	}, 120*time.Second); err != nil {
		return
	}

	outp, err = cmd.Exec("run", "--image", "localhost:5000/cells:v3", "./services/test/routes/routes-client")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find routes-client in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}
		if !statusRunning("routes-client", "latest", outp) {
			return outp, errors.New("Can't find routes-client in runtime")
		}
		return nil, nil
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Check logs", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "routes-client")
		if err != nil {
			return nil, err
		}
		if !strings.Contains(string(outp), "Client completed ok") {
			return outp, fmt.Errorf("Client did not complete ok")
		}
		return nil, nil
	}, 120*time.Second); err != nil {
		return
	}
}
