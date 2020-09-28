// +build integration

package test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRPC(t *testing.T) {
	TrySuite(t, testRPC, retryCount)
}

func testRPC(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	outp, err := cmd.Exec("run", "./service/rpc/rpc-server")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find rpc-server", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "rpc") {
			return outp, errors.New("Can't find rpc service in registry")
		}
		return nil, nil
	}, 90*time.Second); err != nil {
		return
	}

	outp, err = cmd.Exec("run", "./service/rpc/rpc-client")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Check logs", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "rpc-client")
		if err != nil {
			return nil, err
		}
		if !strings.Contains(string(outp), "Client completed ok") {
			return outp, fmt.Errorf("Client did not complete ok")
		}
		return nil, nil
	}, 60*time.Second); err != nil {
		return
	}
}
