// +build integration

package test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestEnvOverrides(t *testing.T) {
	trySuite(t, testEnvOverrides, retryCount)
}

func testEnvOverrides(t *t) {
	outp, err := exec.Command("micro", "-env=platform", "env").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		return
	}
	if !strings.Contains(string(outp), "* platform") {
		t.Fatal("Env platform is not selected")
		return
	}

	outp, err = exec.Command("micro", "-e=platform", "env").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		return
	}
	if !strings.Contains(string(outp), "* platform") {
		t.Fatal("Env platform is not selected")
		return
	}
}
