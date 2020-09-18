// +build integration

package test

import (
	"errors"
	"fmt"
	"os"
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

	runTarget := "./service/rpc/rpc-server"
	branch := "latest"
	if os.Getenv("MICRO_IS_KIND_TEST") == "true" {
		if ref := os.Getenv("GITHUB_REF"); len(ref) > 0 {
			branch = strings.TrimPrefix(ref, "refs/heads/")
		} else {
			branch = "master"
		}
		runTarget = "github.com/micro/micro/test/service/rpc/rpc-server@" + branch
		t.Logf("Running service from the %v branch of micro", branch)
	}

	outp, err := cmd.Exec("run", runTarget)
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

	runTarget = "./service/rpc/rpc-client"
	branch = "latest"
	if os.Getenv("MICRO_IS_KIND_TEST") == "true" {
		if ref := os.Getenv("GITHUB_REF"); len(ref) > 0 {
			branch = strings.TrimPrefix(ref, "refs/heads/")
		} else {
			branch = "master"
		}
		runTarget = "github.com/micro/micro/test/service/rpc/rpc-client@" + branch
		t.Logf("Running service from the %v branch of micro", branch)
	}

	outp, err = cmd.Exec("run", runTarget)
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
