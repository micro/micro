// +build integration

package test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestEnvBasic(t *testing.T) {
	trySuite(t, testEnvOverrides, retryCount)
}

func testEnvBasic(t *t) {
	outp, err := exec.Command("micro", "-env=platform", "env").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		return
	}
	if !strings.Contains(string(outp), "platform") ||
		!strings.Contains(string(outp), "server") ||
		!strings.Contains(string(outp), "local") {
		t.Fatal("Env output lacks local, server, or platform")
		return
	}

}

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
