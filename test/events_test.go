// +build integration

package test

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestEventsStream(t *testing.T) {
	TrySuite(t, testEventsStream, RetryCount)
}

func testEventsStream(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin(), WithNamespace("micro"))
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// Temp fix to support k8s tests until we have file upload to remote server
	var branch string
	if ref := os.Getenv("GITHUB_REF"); len(ref) > 0 {
		branch = strings.TrimPrefix(ref, "refs/heads/")
	} else {
		branch = "master"
	}

	t.Logf("Running service from the %v branch of micro", branch)
	if outp, err := cmd.Exec("run", "github.com/micro/micro/test/service/stream@"+branch); err != nil {
		t.Fatalf("Error running service: %v, %v", err, string(outp))
		return
	}

	if err := Try("Check logs for success", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "-n", "200", "stream")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "Published event ok") {
			return outp, errors.New("Published event log not found")
		}
		if !strings.Contains(string(outp), "Recieved event ok") {
			return outp, errors.New("Recieved event log not found")
		}
		return outp, nil
	}, 180*time.Second); err != nil {
		return
	}
}
