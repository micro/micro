// +build integration

package test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestFileUpload(t *testing.T) {
	trySuite(t, testFileUpload, retryCount)
}

func testFileUpload(t *t) {
	t.Parallel()
	serv := newServer(t, options{
		auth: "jwt",
	})
	defer serv.close()
	if err := serv.launch(); err != nil {
		return
	}

	login(serv, t, "default", "password")

	outp, err := exec.Command("micro", serv.envFlag(), "run", "./example-service").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	if err := try("Test store for existence of the file", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "store", "list", "--table", "server").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "files/micro/example-service.tar.gz") {
			return outp, fmt.Errorf("Output should contain example service")
		}
		return outp, nil
	}, 15*time.Second); err != nil {
		return
	}
}
