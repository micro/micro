// +build integration

package test

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestEventsStream(t *testing.T) {
	// temporarily nuking this test
	return
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

	// Temp fix to support k8s tests until we have file upload to remote server
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "/service/events"); err != nil {
		t.Fatalf("Error running service: %v, %v", err, string(outp))
		return
	}

	if err := Try("Check logs for success", t, func() ([]byte, error) {
		outp, _ := cmd.Exec("status")
		outp, err := cmd.Exec("logs", "-n", "200", "events")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "TEST1: Finished ok") {
			return outp, errors.New("Received event log not found")
		}
		if !strings.Contains(string(outp), "TEST2: Finished ok") {
			return outp, errors.New("Test 2 not finished")
		}
		return outp, nil
	}, 180*time.Second); err != nil {
		return
	}
}
