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
	TrySuite(t, testEventsStream, retryCount)
}

func testEventsStream(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin(), WithNamespace("micro"))
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	var err error
	var outp []byte
	// Temp fix to support k8s tests until we have file upload to remote server
	if ref := os.Getenv("GITHUB_REF"); len(ref) > 0 {
		ref = strings.TrimPrefix(ref, "refs/heads/")
		t.Logf("Running service from the %v branch of micro", ref)
		outp, err = cmd.Exec("run", "--image", "localhost:5000/cells:micro", "github.com/micro/micro/test/service/stream@"+ref)
	} else {
		outp, err = cmd.Exec("run", "--image", "localhost:5000/cells:micro", "./service/stream")
	}
	if err != nil {
		t.Fatalf("Error running service: %v, %v", err, string(outp))
		return
	}

	if err := Try("Check logs for success", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "-n", "200", "stream")
		if err != nil {
			return outp, err
		}
		t.Logf(string(outp))
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
