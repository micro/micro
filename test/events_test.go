// +build integration

package test

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestEventsStream(t *testing.T) {
	TrySuite(t, testEventsStream, retryCount)
}

func testEventsStream(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "./service/stream")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Check logs for success", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "test/stream")
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
	}, 90*time.Second); err != nil {
		return
	}
}
