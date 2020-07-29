// +build integration

package test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestEnvBasic(t *testing.T) {
	TrySuite(t, testEnvOverrides, retryCount)
}

func testEnvBasic(t *T) {
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
	TrySuite(t, testEnvOverrides, retryCount)
}

func testEnvOverrides(t *T) {
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

func TestEnvOps(t *testing.T) {
	TrySuite(t, testEnvOps, retryCount)
}

func testEnvOps(t *T) {
	// add an env
	_, err := exec.Command("micro", "env", "add", "fooTestEnvOps", "127.0.0.1:8081").CombinedOutput()
	if err != nil {
		t.Fatalf("Error running micro env %s", err)
		return
	}
	outp, err := exec.Command("micro", "env").CombinedOutput()
	if err != nil {
		t.Fatalf("Error running micro env %s", err)
		return
	}
	if !strings.Contains(string(outp), "fooTestEnvOps") {
		t.Fatalf("Cannot find expected environment. Output %s", outp)
		return
	}

	// can we actually set it correctly
	_, err = exec.Command("micro", "env", "set", "fooTestEnvOps").CombinedOutput()
	if err != nil {
		t.Fatalf("Error running micro env %s", err)
		return
	}
	outp, err = exec.Command("micro", "env").CombinedOutput()
	if err != nil {
		t.Fatalf("Error running micro env %s", err)
		return
	}
	if !strings.Contains(string(outp), "* fooTestEnvOps") {
		t.Fatalf("Environment not set. Output %s", outp)
		return
	}

	_, err = exec.Command("micro", "env", "set", "local").CombinedOutput()
	if err != nil {
		t.Fatalf("Error running micro env %s", err)
		return
	}

	// we should be able to delete it too
	_, err = exec.Command("micro", "env", "del", "fooTestEnvOps").CombinedOutput()
	if err != nil {
		t.Fatalf("Error running micro env %s", err)
		return
	}

	outp, err = exec.Command("micro", "env").CombinedOutput()
	if err != nil {
		t.Fatalf("Error running micro env %s", err)
		return
	}
	if strings.Contains(string(outp), "fooTestEnvOps") {
		t.Fatalf("Found unexpected environment. Output %s", outp)
		return
	}

}
